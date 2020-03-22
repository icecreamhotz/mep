package utils

import (
	"strconv"
	"time"

	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"
)

type NotFound struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

type ServerError struct {
	Success bool     `json:"success"`
	Message string   `json:"message"`
	Data    []string `json:"data"`
}

type Token struct {
	Success      bool      `json:"success"`
	Message      string    `json:"message"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredAt    time.Time `json:"expired_at"`
}

type UnProcessEntity struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Errors  []map[string]string `json:"errors"`
}

type Message struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func ResponseObject(obj interface{}) interface{} {
	return obj
}

func ResponseMessage(message string) Message {
	return Message{
		Success: false,
		Message: message,
	}
}

func ReponseNotFound(message string) NotFound {
	return NotFound{
		Success: false,
		Message: message,
		Data:    []string{},
	}
}

func ResponseServerError(message string) ServerError {
	return ServerError{
		Success: false,
		Message: message,
		Data:    []string{},
	}
}

func ResponseToken(message, access_token, refresh_token string, expired_at time.Time) Token {
	return Token{
		Success:      true,
		Message:      message,
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		ExpiredAt:    expired_at,
	}
}

func ResponseErrorFields(errors []map[string]string) UnProcessEntity {
	return UnProcessEntity{
		Success: false,
		Message: "Invalid Fields",
		Errors:  errors,
	}
}

func ResponseErrorValidation(translator ut.Translator, err error) UnProcessEntity {
	if _, ok := err.(*strconv.NumError); ok {
		return ResponseErrorFields([]map[string]string{{
			"done": "Done must be accept boolean value only",
		}})
	}
	errs := err.(validator.ValidationErrors)
	invalidFields := make([]map[string]string, 0)
	for _, e := range errs {
		errors := map[string]string{}
		errors[ToSnakeCase(e.Field())] = e.Translate(translator)
		invalidFields = append(invalidFields, errors)
	}
	return ResponseErrorFields(invalidFields)
}
