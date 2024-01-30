package module

import (
	"github.com/GoPOS-id/gopos-api/module/auth"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	auth.Register(app)
}
