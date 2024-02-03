package utils

import "github.com/gofiber/fiber/v2"

type DataResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func SendDataResponse(c *fiber.Ctx, message string, data interface{}, code int) error {
	json := DataResponse{
		Code:    getStatus(code),
		Message: message,
		Data:    data,
	}
	return c.Status(getStatus(code)).JSON(json)
}

func SendResponse(c *fiber.Ctx, message string, code int) error {
	json := Response{
		Code:    getStatus(code),
		Message: message,
	}
	return c.Status(getStatus(code)).JSON(json)
}

func getStatus(code int) int {
	if code == 0 {
		return 200
	}
	return code
}
