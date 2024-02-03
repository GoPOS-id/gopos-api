package auth

import (
	"github.com/GoPOS-id/gopos-api/utils"
	"github.com/gofiber/fiber/v2"
)

// @Tags Authentication
// @Summary Login
// @Description Authenticates a user and returns an access token.
// @ID loginUser
// @Accept json
// @Produce json
// @Param body body inAuthDtos true "User credentials for login"
// @Success 200 {object} utils.DataResponse{data=[]outAuthDtos} "User login successful"
// @Failure 400 {object} utils.Response "Invalid input data"
// @Router /auth [post]
// loginHandler handles the user login request.
// It parses the request body to retrieve authentication data,
// validates the input data, performs user login, and sends a response accordingly.
func loginHandler(c *fiber.Ctx) error {
	// Parse the request body to retrieve authentication data
	dtos := new(inAuthDtos)
	if err := c.BodyParser(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 409)
	}

	// Validate the input data for login
	if err := handleLoginValidator(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 400)
	}

	// Perform user login
	user, err := dtos.login()
	if err != nil {
		return handleLoginError(c, err)
	}

	// Send a successful login response
	return utils.SendDataResponse(c, "Login successfully", user, 200)
}

// @Tags Authentication
// @Summary Logout
// @Description Logs out the user by invalidating the provided access token.
// @ID logoutUser
// @Produce json
// @Param Authorization header string true "Bearer <access_token>" default("Bearer ")
// @Success 200 {object} utils.Response "Logout Successfully"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal Server Error"
// @Router /auth [delete]
// logoutHandler handles the user logout request.
// It extracts the JWT token from the authorization header,
// initiates the logout process, and sends a response accordingly.
func logoutHandler(c *fiber.Ctx) error {
	// Extract the JWT token from the Authorization header
	authorization := c.Get("Authorization")
	token := handleToken(authorization)

	// Create an inAuthDtos object with the extracted token
	dtos := outAuthDtos{
		Token: token,
	}

	// Perform user logout
	if err := dtos.logout(); err != nil {
		return utils.SendResponse(c, err.Error(), 500)
	}

	// Send a successful logout response
	return utils.SendResponse(c, "Logout Successfully", 200)
}
