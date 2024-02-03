package server

import (
	"errors"

	"github.com/GoPOS-id/gopos-api/utils"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func fiberConfigrutaion() *fiber.App {
	engine := html.New("./views", ".html")
	return fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Views:       engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := "Internal Server Error"

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
				if code == 404 {
					msg = "Not Found"
				} else if code == 500 {
					msg = "Internal Server Error"
				} else {
					msg = e.Message
				}
			}
			err = utils.SendResponse(c, msg, code)

			if err != nil {
				return utils.SendResponse(c, "Internal Server Error", 500)
			}

			return nil
		},
	})
}
