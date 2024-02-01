package auth

import (
	"strings"

	"github.com/GoPOS-id/gopos-api/utils"
	"github.com/gofiber/fiber/v2"
)

// ! LOGOUT HANDLERS
func logoutHandler(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	jwtsplit := strings.Split(token, ".")

	dtos := inAuthDtos{
		Token: jwtsplit[2],
	}

	if err := dtos.Logout(); err != nil {
		resp := utils.DataResponse{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		}
		return resp.SendMessageJSON(c)
	}

	resp := utils.DataResponse{
		Code:    200,
		Message: "Logout successfully",
	}
	return resp.SendMessageJSON(c)
}
