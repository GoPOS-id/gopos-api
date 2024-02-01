package server

import (
	"errors"

	"github.com/GoPOS-id/gopos-api/utils"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func fiberConfigrutaion() *fiber.App {
	return fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			msg := "Internal Server Error"

			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
				if code == fiber.StatusNotFound {
					msg = "Not Found"
				} else if code == fiber.StatusInternalServerError {
					msg = "Internal Server Error"
				} else {
					msg = e.Message
				}
			}

			resp := utils.DataResponse{
				Code:    code,
				Message: msg,
			}
			err = resp.SendMessageJSON(c)

			if err != nil {
				resp := utils.DataResponse{
					Code:    500,
					Message: "Internal Server Error",
				}
				return resp.SendMessageJSON(c)
			}

			return nil
		},
	})
}
