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
	UserHandler controllers.UserHandler
	AuthHandler controllers.AuthHandler
}

func NewValidateTranslation() ut.Translator {
	var trans ut.Translator
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		en := en.New()
		uni := ut.New(en, en)
		trans, _ = uni.GetTranslator("en")
		en_translations.RegisterDefaultTranslations(v, trans)
	}

	return trans
}

func NewApplication(userHandler controllers.UserHandler, authHandler controllers.AuthHandler) App {
	return App{
		UserHandler: userHandler,
		AuthHandler: authHandler,
	}
}
