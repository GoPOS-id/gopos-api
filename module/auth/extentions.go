package auth

import (
	"strings"

	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func handleLoginValidator(dtos *inAuthDtos) error {
	return validation.ValidateStruct(dtos,
		validation.Field(&dtos.Username, validation.Required, validation.Length(8, 25)),
		validation.Field(&dtos.Password, validation.Required, validation.Length(8, 25)),
	)
}

func handleLoginError(c *fiber.Ctx, err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return utils.SendResponse(c, "Invalid username or password", 400)
	default:
		return utils.SendResponse(c, err.Error(), 500)
	}
}

func handleToken(authorization string) string {
	token := strings.TrimPrefix(authorization, "Bearer ")
	jwtsplit := strings.Split(token, ".")

	return jwtsplit[2]
}

func mapSessionDto(userId int64, token string) model.Session {
	return model.Session{
		UserId: userId,
		Token:  strings.Split(token, ".")[2],
	}
}

func mapOutDto(userId int64, username string, token string) outAuthDtos {
	return outAuthDtos{
		Id:       userId,
		Username: username,
		Token:    token,
	}
}
