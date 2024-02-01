package auth

import (
	"strings"

	"github.com/GoPOS-id/gopos-api/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ! LOGIN HANDLERS
func loginHandler(c *fiber.Ctx) error {
	dtos := new(inAuthDtos)
	if err := c.BodyParser(dtos); err != nil {
		resp := utils.DataResponse{
			Code:    fiber.StatusConflict,
			Message: err.Error(),
		}
		return resp.SendMessageJSON(c)
	}

	validate := validation.ValidateStruct(dtos,
		validation.Field(&dtos.Username, validation.Required, validation.Length(8, 25)),
		validation.Field(&dtos.Password, validation.Required, validation.Length(8, 25)),
	)

	if validate != nil {
		resp := utils.DataResponse{
			Code:    fiber.StatusBadRequest,
			Message: validate.Error(),
		}
		return resp.SendMessageJSON(c)
	}

	user, err := dtos.Login()

	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			resp := utils.DataResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Username or Password invalid",
			}
			return resp.SendMessageJSON(c)
		default:
			resp := utils.DataResponse{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			}
			return resp.SendMessageJSON(c)
		}
	}
	resp := utils.DataResponse{
		Message: "Login successfully",
		Data:    user,
	}

	return resp.SendDataJSON(c)
}

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
