package lib

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/caffeines/filepile/config"
	"github.com/dgrijalva/jwt-go"
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
