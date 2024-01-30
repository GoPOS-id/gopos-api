package utils

import "github.com/gofiber/fiber/v2"

type DataResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (resp *DataResponse) SendDataJSON(c *fiber.Ctx) error {
	return c.Status(200).JSON(resp)
}

func (resp *DataResponse) SendMessageJSON(c *fiber.Ctx) error {
	return c.Status(resp.Code).JSON(fiber.Map{
		"code":    resp.Code,
		"message": resp.Message,
	})
}
