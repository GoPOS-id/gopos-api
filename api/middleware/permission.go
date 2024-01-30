package middleware

import (
	"github.com/GoPOS-id/gopos-api/constant"
	"github.com/GoPOS-id/gopos-api/utils"
	"github.com/gofiber/fiber/v2"
)

var resp = utils.DataResponse{
	Code:    fiber.StatusBadRequest,
	Message: "You didn't have permission",
}

func OperatorOnly(c *fiber.Ctx) error {
	locals := c.Locals("user").(OutAuthDtos)
	if locals.Role != constant.Role_Operator {
		return resp.SendMessageJSON(c)
	}
	return c.Next()
}

func AdministratorOnly(c *fiber.Ctx) error {
	locals := c.Locals("user").(OutAuthDtos)
	op := constant.Role_Operator
	ad := constant.Role_Adminstrator
	if locals.Role != ad || locals.Role != op {
		return resp.SendMessageJSON(c)
	}
	return c.Next()
}
