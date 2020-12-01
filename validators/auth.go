package validators

import (
	"fmt"

	"github.com/caffeines/sharehub/lib"
	"github.com/caffeines/sharehub/models"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/labstack/echo/v4"
)

var uni *ut.UniversalTranslator

func ValidateRegister(ctx echo.Context) (*models.User, error) {
	ur := struct {
		Name     string `json:"name,omitempty" validate:"required,min=2,max=50"`
		Username string `json:"username,omitempty" validate:"required,min=5,max=26"`
		Email    string `json:"email,omitempty" validate:"required,email"`
		Password string `json:"password,omitempty" validate:"required,min=6,max=26"`
	}{}
	en := en.New()
	uni = ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	if err := ctx.Bind(&ur); err != nil {
		return nil, err
	}
	fmt.Printf("%+v", ur)
	v := validator.New()
	en_translations.RegisterDefaultTranslations(v, trans)
	ve := lib.ValidationError{}
	if err := v.Struct(ur); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			fmt.Println(e.Translate(trans))
			ve.Add(e.Field(), e.Translate(trans))
		}
		fmt.Println(ve.Error())
		return nil, err
	}
	return nil, nil
}
