package api

import (
	"log"
	"net/http"
	"time"

	"github.com/caffeines/filepile/app"
	"github.com/caffeines/filepile/constants"
	"github.com/caffeines/filepile/constants/errors"
	"github.com/caffeines/filepile/data"
	"github.com/caffeines/filepile/lib"
	"github.com/caffeines/filepile/middlewares"
	"github.com/caffeines/filepile/models"
	"github.com/caffeines/filepile/validators"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterAuthRoutes registers authintacation routes
func RegisterAuthRoutes(endpoint *echo.Group) {
	endpoint.POST("/register/", register)
	endpoint.POST("/login/", login)
	endpoint.PATCH("/refresh-token/", refreshToken)
	endpoint.PATCH("/logout/", logout, middlewares.JWTAuth())
}

func login(ctx echo.Context) error {
	resp := lib.Response{}
	body, err := validators.ValidateLogin(ctx)
	if err != nil {
		log.Println(err)
		resp.Title = "Invalid login request data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.InvalidRegisterData
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	db := app.GetDB()
	userRepo := data.NewUserRepo()
	user, err := userRepo.FindUserByEmail(db, body.Email)
	if err != nil {
		log.Println(err)
		if lib.IsDocumentNotFoundError(err) {
			resp.Title = "User not found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.UserNotFound
			resp.Errors = err
			return resp.ServerJSON(ctx)
		}
		resp.Title = "User login failed"
		resp.Status = http.StatusUnauthorized
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	if ok := lib.CheckPasswordHash(body.Password, user.Password); !ok {
		resp.Title = "Email or password incorrect"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.InvalidLoginCredential
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	signedToken, err := lib.BuildJWTToken(user.Username, constants.USER_SCOPE, user.ID.Hex())
	if err != nil {
		log.Println(err)

		resp.Title = "Failed to sign auth token"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.UserLoginFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	sess := &models.Session{
		ID:           primitive.NewObjectID(),
		UserID:       user.ID,
		RefreshToken: lib.NewRefresToken(),
		AccessToken:  signedToken,
		CreatedAt:    time.Now().UTC(),
		ExpiresOn:    time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	sessRepo := data.NewSessionRepo()
	if err = sessRepo.CreateSession(db, sess); err != nil {
		log.Println(err)
		resp.Title = "User login failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	result := map[string]interface{}{
		"access_token":  sess.AccessToken,
		"refresh_token": sess.RefreshToken,
		"expire_on":     sess.ExpiresOn,
		"permission":    constants.USER_SCOPE,
	}
	resp.Status = http.StatusOK
	resp.Data = result
	return resp.ServerJSON(ctx)
}

func register(ctx echo.Context) error {
	resp := lib.Response{}
	user, err := validators.ValidateRegister(ctx)
	if err != nil {
		log.Println(err)
		resp.Title = "Invalid request data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.InvalidRegisterData
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	hash, err := lib.HashPassword(user.Password)
	if err != nil {
		log.Println(err)
		resp.Title = "User registration failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	user.Password = hash
	userRepo := data.NewUserRepo()
	db := app.GetDB()
	_, err = userRepo.Register(db, user)

	if err != nil {
		log.Println(err)
		resp.Title = "User registration failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	resp.Title = "User registration successful"
	resp.Status = http.StatusOK
	resp.Data = user

	return resp.ServerJSON(ctx)
}

func refreshToken(ctx echo.Context) error {
	resp := lib.Response{}
	token, err := lib.ParseRefreshToken(ctx)
	if err != nil {
		resp.Title = "Token parsing failed"
		resp.Errors = err
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = errors.UserSignUpDataInvalid
		return resp.ServerJSON(ctx)
	}
	db := app.GetDB()
	sessionRepo := data.NewSessionRepo()
	claims, _, err := lib.ExtractAndValidateToken(ctx)
	if err != nil {
		resp.Title = "Bearer token not found or expired"
		resp.Status = http.StatusNotFound
		resp.Code = errors.BearerTokenNotFound
		resp.Errors = lib.NewError(err.Error())
		return resp.ServerJSON(ctx)
	}
	accessToken, err := lib.BuildJWTToken(claims.Username, constants.USER_SCOPE, claims.UserID)
	if err != nil {
		resp.Title = "Failed to sign auth token"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.TokenRefreshFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	sess, err := sessionRepo.UpdateSession(db, token, accessToken)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			resp.Title = "Refresh token not found or expired"
			resp.Status = http.StatusNotFound
			resp.Code = errors.RefreshTokenNotFound
			resp.Errors = lib.NewError(err.Error())
			return resp.ServerJSON(ctx)
		}
		resp.Title = "Token refresh failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	resp.Data = sess
	resp.Status = http.StatusOK
	return resp.ServerJSON(ctx)
}

func logout(ctx echo.Context) error {
	resp := lib.Response{}
	token, err := lib.ParseRefreshToken(ctx)
	if err != nil {
		resp.Title = "Invalid token data"
		resp.Status = http.StatusBadRequest
		resp.Code = errors.BearerTokenNotFound
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	log.Println(token)
	db := app.GetDB()
	sessionRepo := data.NewSessionRepo()
	if err := sessionRepo.Logout(db, token); err != nil {
		if err == mongo.ErrNoDocuments {
			resp.Title = "No session found"
			resp.Status = http.StatusNotFound
			resp.Code = errors.RefreshTokenNotFound
			resp.Errors = lib.NewError(err.Error())
			return resp.ServerJSON(ctx)
		}
		resp.Title = "Logout failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = errors.DatabaseQueryFailed
		resp.Errors = err
		return resp.ServerJSON(ctx)
	}
	resp.Status = http.StatusOK
	resp.Title = "Logout successful"
	return resp.ServerJSON(ctx)
}
