package user

import (
	"github.com/GoPOS-id/gopos-api/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	route := app.Group("/user")
	route.Use(middleware.Auth)
	route.Get("/", getProfileHandler)
	route.Post("/", createUserHandler)
	route.Patch("/", middleware.AdministratorOnly, updateUserHandler)

}
