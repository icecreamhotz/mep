package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	ut "github.com/go-playground/universal-translator"
	"github.com/icecreamhotz/mep-api/models"
	"github.com/icecreamhotz/mep-api/utils"
)

type UserHandler struct {
	Service   models.UserReporer
	Validator ut.Translator
}

func NewUserHandler(repository models.UserReporer, validator ut.Translator) UserHandler {
	return UserHandler{
		Service:   repository,
		Validator: validator,
	}
}

func (handler *UserHandler) BackofficeUserPost(c *gin.Context) {
	var bofRequest models.BackofficeUser
	var err error

	if err = c.ShouldBind(&bofRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ResponseErrorValidation(handler.Validator, err))
		return
	}
	// find username exists
	backofficeUser, err := handler.Service.FindByUsername(bofRequest.Username)
	if err != nil && err != pg.ErrNoRows {
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
		return
	}

	if backofficeUser.ID.String() != "" {
		c.JSON(http.StatusConflict, utils.ResponseErrorFields([]map[string]string{{
			"username": "Duplicate Username.",
		}}))
		return
	}

	// find email exists
	backofficeUser, err = handler.Service.FindByEmail(bofRequest.Email)
	if err != nil && err != pg.ErrNoRows {
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
		return
	}

	if backofficeUser.ID.String() != "" {
		c.JSON(http.StatusConflict, utils.ResponseErrorFields([]map[string]string{{
			"username": "Duplicate Email.",
		}}))
		return
	}

	hashPassword, errHashPassword := utils.HashPassword(bofRequest.Password)
	if errHashPassword != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
	}

	bofRequest.Password = hashPassword

	err = handler.Service.Create(bofRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
	}
	c.JSON(http.StatusCreated, utils.ResponseMessage("Created successful."))
}
