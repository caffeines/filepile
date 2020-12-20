package lib

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/caffeines/filepile/config"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type Claims struct {
	UserID   string `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func BuildJWTToken(username, scope, id string) (string, error) {
	claims := Claims{
		UserID:   id,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Audience:  scope,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.GetJWT().TTL)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetJWT().Secret))
}

func NewRefresToken() string {
	token := fmt.Sprintf("%d_%s", time.Now().Unix(), NewUUID())
	return base64.StdEncoding.EncodeToString([]byte(token))
}
func extractTokenFromHeader(ctx echo.Context) string {
	tokenWithBearer := ctx.Request().Header.Get("Authorization")
	token := strings.Replace(tokenWithBearer, "Bearer", "", -1)
	return strings.TrimSpace(token)
}

func ExtractAndValidateToken(ctx echo.Context) (*Claims, *jwt.Token, error) {
	token := extractTokenFromHeader(ctx)
	if token == "" {
		return nil, nil, NewError("Authorization token not found")
	}
	claims := Claims{}
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(config.GetJWT().Secret), nil
	})
	if err != nil {
		return nil, nil, err
	}
	if !jwtToken.Valid {
		return nil, nil, NewError("Token is invalid")
	}
	return &claims, jwtToken, nil
}

func ParseRefreshToken(ctx echo.Context) (string, error) {
	refresh := ctx.Request().Header.Get("RefreshToken")
	refreshWithToken := strings.Split(refresh, " ")
	if len(refreshWithToken) != 2 {
		return "", NewError("Refresh token not found")
	}
	return refreshWithToken[1], nil
}
