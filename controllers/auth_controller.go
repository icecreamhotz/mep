package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v9"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis/v7"
	"github.com/icecreamhotz/mep-api/configs"
	"github.com/icecreamhotz/mep-api/models"
	"github.com/icecreamhotz/mep-api/utils"
)

type AuthHandler struct {
	Service   models.UserReporer
	Validator ut.Translator
}

func NewAuthHandler(repository models.UserReporer, validator ut.Translator) AuthHandler {
	return AuthHandler{
		Service:   repository,
		Validator: validator,
	}
}

func (handler *AuthHandler) AuthLoginPost(c *gin.Context) {
	var credential struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var err error

	if err = c.ShouldBind(&credential); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ResponseErrorValidation(handler.Validator, err))
		return
	}

	backofficeUser, err := handler.Service.FindByUsername(credential.Username)
	if err != nil {
		if err == pg.ErrNoRows {
			c.JSON(http.StatusUnauthorized, utils.ResponseMessage("Unauthorized."))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
		return
	}

	ok := utils.CheckPasswordHash(credential.Password, backofficeUser.Password)
	if !ok {
		c.JSON(http.StatusUnauthorized, utils.ResponseMessage("Unauthorized."))
		return
	}

	accessToken, refreshToken, expirationTimeAccessToken, err := utils.SetAccessTokenAndRefreshToken(
		backofficeUser.ID,
		backofficeUser.Name,
		backofficeUser.Lastname,
		backofficeUser.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
	}

	c.JSON(http.StatusOK, utils.ResponseToken("Login successful.",
		accessToken,
		refreshToken,
		expirationTimeAccessToken))
}

func (handler *AuthHandler) AuthRefreshTokenPost(c *gin.Context) {
	var token struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	var err error

	if err = c.ShouldBind(&token); err != nil {
		c.JSON(http.StatusUnprocessableEntity, utils.ResponseErrorValidation(handler.Validator, err))
		return
	}

	client := configs.ConnectCacheDatabase()

	userId, err := client.HGet("refresh_token:"+token.RefreshToken, "user_id").Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusUnauthorized, utils.ResponseMessage("Unauthorized."))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
		return
	}

	backofficeUser, err := handler.Service.GetById(userId)
	if err != nil {
		if err == pg.ErrNoRows {
			c.JSON(http.StatusUnauthorized, utils.ResponseMessage("Unauthorized."))
			return
		}
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
		return
	}

	accessToken, refreshToken, expirationTimeAccessToken, err := utils.SetAccessTokenAndRefreshToken(
		backofficeUser.ID,
		backofficeUser.Name,
		backofficeUser.Lastname,
		backofficeUser.Role,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
	}

	err = client.Del("refresh_token:" + token.RefreshToken).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
	}

	c.JSON(http.StatusOK, utils.ResponseToken("Refresh token success.",
		accessToken,
		refreshToken,
		expirationTimeAccessToken))
}

func (handler *AuthHandler) AuthPayloadGet(c *gin.Context) {
	c.JSON(http.StatusOK, utils.ResponseObject(gin.H{
		"id":       c.MustGet("ID"),
		"name":     c.MustGet("Name"),
		"lastname": c.MustGet("Lastname"),
		"role":     c.MustGet("Role"),
	}))
}
