package middlewares

import (
	"net/http"

	"github.com/caffeines/filepile/constants"
	"github.com/caffeines/filepile/constants/errors"
	"github.com/caffeines/filepile/lib"
	"github.com/labstack/echo/v4"
)

func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			resp := lib.Response{}

			claims, _, err := lib.ExtractAndValidateToken(ctx)
			if err != nil {
				resp.Status = http.StatusUnauthorized
				resp.Code = errors.InvalidAuthorizationToken
				resp.Title = "Unauthorized request"
				resp.Errors = err
				return resp.ServerJSON(ctx)
			}
			ctx.Set(constants.USER_ID, claims.UserID)
			ctx.Set(constants.USER_SCOPE, claims.Audience)
			ctx.Set(constants.USERNAME, claims.Username)
			return next(ctx)
		}
	}
}
