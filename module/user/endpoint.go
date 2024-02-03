package user

import (
	"github.com/GoPOS-id/gopos-api/api/middleware"
	"github.com/gofiber/fiber/v2"
)

// Register sets up user-related routes for the given Fiber app.
func Register(app *fiber.App) {
	route := app.Group("/user")                                         // Create a route group for user-related endpoints under the "/user" path.
	route.Use(middleware.Auth)                                          // Use authentication middleware for all endpoints in the "/user" group.
	route.Get("/", getProfileHandler)                                   // GET endpoint to retrieve the profile of the authenticated user.
	route.Get("/all", middleware.AdministratorOnly, getAllUsersHandler) // GET endpoint to retrieve information about all users. Restricted to administrators only.
	route.Post("/", middleware.AdministratorOnly, createUserHandler)    // POST endpoint to create a new user. Restricted to administrators only.
	route.Patch("/", middleware.AdministratorOnly, updateUserHandler)   // PATCH endpoint to update user information. Restricted to administrators only.
	route.Delete("/", middleware.AdministratorOnly)                     // DELETE endpoint to delete a user. Restricted to administrators only.
}
