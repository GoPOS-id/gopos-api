package auth

import (
	"github.com/GoPOS-id/gopos-api/api/middleware"
	"github.com/gofiber/fiber/v2"
)

func Register(app *fiber.App) {
	route := app.Group("/auth")
	route.Post("/", loginHandler)
	route.Delete("/", middleware.Auth, logoutHandler)
}
