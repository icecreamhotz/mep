package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/icecreamhotz/mep-api/configs"
	"github.com/icecreamhotz/mep-api/routes"
	"github.com/spf13/viper"
)

func main() {
	if err := configs.InitConfigs(); err != nil {
		log.Fatal("%s", err.Error())
	}

	app, err := InitialApplication()
	if err != nil {
		log.Fatal("%s", err.Error())
	}

	router := gin.Default()
	v1 := router.Group("/api/v1")
	routes.RouterBackofficeUser(v1, app.UserHandler)
	routes.RouterAuthenticate(v1, app.AuthHandler)
	routes.RouterTodoLists(v1, app.TodoListHandler)

	router.Run(":" + viper.GetString("port"))
}
