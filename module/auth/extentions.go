package auth

import (
	"strings"

	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// handleLoginValidator performs validation on login data.
// It checks if the username and password meet the required length criteria.
// Returns an error if validation fails or nil otherwise.
func handleLoginValidator(dtos *inAuthDtos) error {
	return validation.ValidateStruct(dtos,
		validation.Field(&dtos.Username, validation.Required, validation.Length(8, 25)),
		validation.Field(&dtos.Password, validation.Required, validation.Length(8, 25)),
	)
}

// handleLoginError handles login-related errors and sends an appropriate response to the client.
// It checks for a specific type of error (gorm.ErrRecordNotFound) and sends a 400 response for invalid username or password.
// For other errors, it sends a generic 500 response with the error message.
func handleLoginError(c *fiber.Ctx, err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return utils.SendResponse(c, "Invalid username or password", 400)
	default:
		return utils.SendResponse(c, err.Error(), 500)
	}
}

// handleToken extracts the JWT token from the 'Authorization' header.
// It removes the 'Bearer ' prefix and splits the token into its three parts.
// Returns the third part of the JWT token.
func handleToken(authorization string) string {
	token := strings.TrimPrefix(authorization, "Bearer ")
	jwtsplit := strings.Split(token, ".")

	return jwtsplit[2]
}

// mapSessionDto creates a 'Session' model from the provided user ID and JWT token.
// Returns the mapped 'Session' model.
func mapSessionDto(userId int64, token string) model.Session {
	return model.Session{
		UserId: userId,
		Token:  strings.Split(token, ".")[2],
	}
}

// mapOutDto creates an 'outAuthDtos' model from the provided user ID, username, and JWT token.
// Returns the mapped 'outAuthDtos' model.
func mapOutDto(userId int64, username string, token string) outAuthDtos {
	return outAuthDtos{
		Id:       userId,
		Username: username,
		Token:    token,
	}
}
