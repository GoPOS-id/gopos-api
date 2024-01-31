package module

import (
	"github.com/GoPOS-id/gopos-api/app/module/auth"
	"github.com/GoPOS-id/gopos-api/app/module/user"
	"github.com/gofiber/fiber/v2"
)

func Init(app *fiber.App) {
	auth.Register(app)
	user.Register(app)
}
