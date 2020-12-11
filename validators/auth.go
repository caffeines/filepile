package validators

import (
	"time"

	"github.com/caffeines/filepile/lib"
	"github.com/caffeines/filepile/models"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateRegister(ctx echo.Context) (*models.User, error) {
	usr := struct {
		Name     string `json:"name,omitempty" validate:"required,min=2,max=50"`
		Username string `json:"username,omitempty" validate:"required,min=5,max=26"`
		Email    string `json:"email,omitempty" validate:"required,email"`
		Password string `json:"password,omitempty" validate:"required,min=6,max=26"`
	}{}

	if err := ctx.Bind(&usr); err != nil {
		return nil, err
	}
	if err := lib.GetValidationError(usr); err != nil {
		return nil, err
	}
	user := &models.User{
		ID:        primitive.NewObjectID(),
		Name:      usr.Name,
		Username:  usr.Username,
		Email:     usr.Email,
		Password:  usr.Password,
		CreatedAt: time.Now().UTC(),
		UpdateAt:  time.Now().UTC(),
	}
	return user, nil
}

// ReqLogin holds login request data
type ReqLogin struct {
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=6,max=26"`
}

// ValidateLogin returns request body or error
func ValidateLogin(ctx echo.Context) (*ReqLogin, error) {
	body := ReqLogin{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := lib.GetValidationError(body); err != nil {
		return nil, err
	}
	return &body, nil
}
