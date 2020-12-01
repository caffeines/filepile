package validators

import (
	"github.com/caffeines/sharehub/lib"
	"github.com/caffeines/sharehub/models"
	"github.com/labstack/echo/v4"
)

func ValidateRegister(ctx echo.Context) (*models.User, error) {
	ur := struct {
		Name     string `json:"name,omitempty" validate:"required,min=2,max=50"`
		Username string `json:"username,omitempty" validate:"required,min=5,max=26"`
		Email    string `json:"email,omitempty" validate:"required,email"`
		Password string `json:"password,omitempty" validate:"required,min=6,max=26"`
	}{}

	if err := ctx.Bind(&ur); err != nil {
		return nil, err
	}
	if err := lib.GetValidationError(ur); err != nil {
		return nil, err
	}
	return nil, nil
}
