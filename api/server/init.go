package server

import (
	"github.com/GoPOS-id/gopos-api/config"
	"github.com/GoPOS-id/gopos-api/database"
	"github.com/GoPOS-id/gopos-api/module"
	"github.com/gofiber/fiber/v2"
)

func Init() {
	database.Init()
	app := fiberConfigrutaion()
	hostname := config.APP_HOST + ":" + config.APP_PORT
	app.Get("/docs", func(c *fiber.Ctx) error {
		return c.Render("docs", fiber.Map{"Title": "GoPOS Api Documentations"})
	})
	app.Get("/docs/json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
	module.Init(app)
	app.Listen(hostname)
}
