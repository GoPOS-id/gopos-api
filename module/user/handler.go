package user

import (
	"github.com/GoPOS-id/gopos-api/api/middleware"
	"github.com/GoPOS-id/gopos-api/constant"
	"github.com/GoPOS-id/gopos-api/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ! GET PROFILE HANDLERS
func getProfileHandler(c *fiber.Ctx) error {
	locals := c.Locals("user").(middleware.OutAuthDtos)
	resp := utils.DataResponse{
		Code:    200,
		Message: "Get Profile Successfully",
		Data:    locals,
	}
	return resp.SendDataJSON(c)
}

// ! CREATE USER HANDLERS
func createUserHandler(c *fiber.Ctx) error {
	dtos := new(inUserDto)
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
		validation.Field(&dtos.Fullname, validation.Required, validation.Length(8, 30)),
		validation.Field(&dtos.Email, validation.Required),
		validation.Field(&dtos.RoleId, validation.Required),
	)

	if validate != nil {
		resp := utils.DataResponse{
			Code:    fiber.StatusBadRequest,
			Message: validate.Error(),
		}
		return resp.SendMessageJSON(c)
	}
	user, err := dtos.Create()
	if err != nil {
		switch err {
		case fiber.ErrBadRequest:
			resp := utils.DataResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Username already exists",
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
		Code:    200,
		Message: "Create user successfully",
		Data:    user,
	}

	return resp.SendDataJSON(c)
}

// ! UPDATE USER HANDLERS
func updateUserHandler(c *fiber.Ctx) error {
	dtos := new(inUserDto)
	locals := c.Locals("user").(middleware.OutAuthDtos)
	if err := c.BodyParser(dtos); err != nil {
		resp := utils.DataResponse{
			Code:    fiber.StatusConflict,
			Message: err.Error(),
		}
		return resp.SendMessageJSON(c)
	}

	validate := validation.ValidateStruct(dtos,
		validation.Field(&dtos.Id, validation.Required),
		validation.Field(&dtos.Username, validation.Required, validation.Length(8, 25)),
		validation.Field(&dtos.Password, validation.Required, validation.Length(8, 25)),
		validation.Field(&dtos.Fullname, validation.Required, validation.Length(8, 30)),
		validation.Field(&dtos.Email, validation.Required),
		validation.Field(&dtos.RoleId, validation.Required),
	)

	if validate != nil {
		resp := utils.DataResponse{
			Code:    fiber.StatusBadRequest,
			Message: validate.Error(),
		}
		return resp.SendMessageJSON(c)
	}

	if dtos.RoleId == 1 || dtos.RoleId == 2 && locals.Role != constant.Role_Operator {
		resp := utils.DataResponse{
			Code:    fiber.StatusBadRequest,
			Message: "Can't change role higher than your currnet role",
		}
		return resp.SendMessageJSON(c)
	}

	user, err := dtos.Update()
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			resp := utils.DataResponse{
				Code:    fiber.StatusBadRequest,
				Message: "Oops!, Can't find the user",
			}
			resp.SendMessageJSON(c)
		default:
			resp := utils.DataResponse{
				Code:    fiber.StatusInternalServerError,
				Message: err.Error(),
			}
			resp.SendMessageJSON(c)
		}
	}
	resp := utils.DataResponse{
		Code:    200,
		Message: "Update changes",
		Data:    user,
	}
	return resp.SendDataJSON(c)
}
