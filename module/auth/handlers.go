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
func loginHandler(c *fiber.Ctx) error {
	dtos := new(inAuthDtos)
	if err := c.BodyParser(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 409)
	}

	if err := handleLoginValidator(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 400)
	}

	user, err := dtos.login()
	if err != nil {
		return handleLoginError(c, err)
	}

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
func logoutHandler(c *fiber.Ctx) error {
	authorization := c.Get("Authorization")
	token := handleToken(authorization)

	dtos := outAuthDtos{
		Token: token,
	}

	if err := dtos.logout(); err != nil {
		return utils.SendResponse(c, err.Error(), 500)
	}

	return utils.SendResponse(c, "Logout Successfully", 200)
}
