package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/icecreamhotz/mep-api/controllers"
	"github.com/icecreamhotz/mep-api/middlewares"
)

func RouterBackofficeUser(router *gin.RouterGroup, handler controllers.UserHandler) {
	backofficeUserRoute := router.Group("/backoffice_user")
	{
		backofficeUserRoute.POST("/", handler.BackofficeUserPost)
	}
}

func RouterAuthenticate(router *gin.RouterGroup, handler controllers.AuthHandler) {
	authRoute := router.Group("/auth")
	{
		authRoute.POST("/login", handler.AuthLoginPost)
		authRoute.POST("/refresh-token", handler.AuthRefreshTokenPost)
		authRoute.GET("/payload", middlewares.Authenticated(), handler.AuthPayloadGet)
	}
}

func RouterTodoLists(router *gin.RouterGroup, handler controllers.TodoListHandler) {
	todoListRoute := router.Group("/todo-list")
	{
		todoListRoute.POST("/", handler.TodoListPost)
		todoListRoute.GET("/", handler.TodoListGet)
		todoListRoute.PATCH("/done/:id", handler.TodoListDonePatch)
	}
}
