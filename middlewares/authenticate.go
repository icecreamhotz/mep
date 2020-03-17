package middlewares

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/icecreamhotz/mep-api/utils"
)

func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenHeader := c.Request.Header.Get("Authorization")

		accessToken, ok := utils.SplitTokenFromHeader(tokenHeader)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseMessage("Unauthorized."))
			return
		}

		// client := configs.ConnectRedis()

		// hasToken, err := client.Exists("access_token:" + accessToken).Result()
		// if err != nil {
		// 	c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
		// 	return
		// }

		// if !hasToken {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseMessage("Unauthorized."))
		// 	return
		// }

		claims := &utils.AccessTokenClaims{}
		claims, parseToken, err := utils.GetUserPayload(accessToken)

		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseMessage("Unauthorized."))
				return
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseMessage("Unauthorized."))
				return
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
				return
			}
		} else if !parseToken.Valid {
			c.AbortWithStatusJSON(http.StatusInternalServerError, utils.ResponseServerError("Something went wrong."))
			return
		}

		c.Set("ID", claims.ID)
		c.Set("Name", claims.Name)
		c.Set("Lastname", claims.Lastname)
		c.Set("Role", claims.Role)

		c.Next()
	}
}
