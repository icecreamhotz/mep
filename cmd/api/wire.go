// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/icecreamhotz/mep-api/configs"
	"github.com/icecreamhotz/mep-api/controllers"
	"github.com/icecreamhotz/mep-api/models"
	"github.com/icecreamhotz/mep-api/server"

	"github.com/google/wire"
)

func InitialApplication() (server.App, error) {
	wire.Build(
		configs.ConfigDatabase,
		configs.NewDatatabase,
		models.NewUserRepository,
		models.NewTodolistRepository,
		server.NewValidateTranslation,
		controllers.NewUserHandler,
		controllers.NewAuthHandler,
		controllers.NewTodoListHandler,
		server.NewApplication,
	)
	return server.App{}, nil
}
