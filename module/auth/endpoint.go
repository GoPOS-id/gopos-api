package auth

import (
	"github.com/GoPOS-id/gopos-api/api/middleware"
	"github.com/gofiber/fiber/v2"
)

// Register sets up authentication-related routes for the given Fiber app.
func Register(app *fiber.App) {
	route := app.Group("/auth")                       // Create a route group for authentication-related endpoints under the "/auth" path.
	route.Post("/", loginHandler)                     // POST endpoint for user login.
	route.Delete("/", middleware.Auth, logoutHandler) // DELETE endpoint for user logout. Requires authentication middleware.
}
