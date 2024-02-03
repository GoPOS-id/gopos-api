package middleware

import (
	"github.com/GoPOS-id/gopos-api/constant"
	"github.com/GoPOS-id/gopos-api/utils"
	"github.com/gofiber/fiber/v2"
)

func OperatorOnly(c *fiber.Ctx) error {
	locals := c.Locals("user").(OutAuthDtos)
	if locals.Role != constant.Role_Operator {
		return utils.SendResponse(c, "You don't have permission", 400)
	}
	return c.Next()
}

func AdministratorOnly(c *fiber.Ctx) error {
	locals := c.Locals("user").(OutAuthDtos)
	op := constant.Role_Operator
	ad := constant.Role_Adminstrator
	if locals.Role == ad || locals.Role == op {
		return c.Next()
	}
	return utils.SendResponse(c, "You don't have permission", 400)
}
