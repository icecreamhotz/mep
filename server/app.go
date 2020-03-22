package server

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/icecreamhotz/mep-api/controllers"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"
)

type App struct {
	UserHandler     controllers.UserHandler
	AuthHandler     controllers.AuthHandler
	TodoListHandler controllers.TodoListHandler
}

func NewValidateTranslation() ut.Translator {
	var trans ut.Translator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		uni := ut.New(en, en)
		trans, _ = uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(v, trans)
		v.RegisterValidation("bool", func(fl validator.FieldLevel) bool {
			if value := fl.Field().String(); value == "" {
				return false
			}
			return true
		})
		v.RegisterTranslation("bool", trans, func(ut ut.Translator) error {
			return ut.Add("bool", "{0} is a required boolean value.", true)
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("bool", fe.Field())
			return t
		})
	}

	return trans
}

func NewApplication(userHandler controllers.UserHandler,
	authHandler controllers.AuthHandler,
	todoListHandler controllers.TodoListHandler) App {
	return App{
		UserHandler:     userHandler,
		AuthHandler:     authHandler,
		TodoListHandler: todoListHandler,
	}
}
