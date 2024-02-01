package user

import (
	"github.com/GoPOS-id/gopos-api/api/middleware"
	"github.com/GoPOS-id/gopos-api/utils"
	"github.com/gofiber/fiber/v2"
)

// ! GET PROFILE HANDLERS
func getProfileHandler(c *fiber.Ctx) error {
	locals := c.Locals("user").(middleware.OutAuthDtos)
	resp := utils.DataResponse{
		Code:    200,
		Message: "Get Profile Successfully",
		Data:    locals,
	}
	return resp.SendDataJSON(c)
}
