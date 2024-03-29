package middleware

import (
	"strings"

	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/config"
	"github.com/GoPOS-id/gopos-api/database"
	"github.com/GoPOS-id/gopos-api/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Auth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {

		return utils.SendResponse(c, "Unauthorized", 401)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	jwtsplit := strings.Split(token, ".")

	if len(jwtsplit) != 3 {
		return utils.SendResponse(c, "Unauthorized", 401)
	}

	verify, err := jwt.ParseWithClaims(token, &utils.JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET_KEY), nil
	})

	if err != nil {
		return utils.SendResponse(c, "Unauthorized", 401)
	}

	if !verify.Valid {
		return utils.SendResponse(c, "Unauthorized", 401)
	}

	user, errors := checkToken(jwtsplit[2])
	if errors != nil {
		switch errors {
		case fiber.ErrUnauthorized:
			return utils.SendResponse(c, "Unauthorized", 401)
		default:
			return utils.SendResponse(c, errors.Error(), 500)
		}
	}
	c.Locals("user", user)
	return c.Next()
}

func checkToken(token string) (OutAuthDtos, error) {
	db := database.DbContext()
	tokenDtos := model.Session{
		Token: token,
	}
	tokenFound := model.Session{}
	if err := db.Model(&tokenDtos).Where("token = ?", tokenDtos.Token).First(&tokenFound).Error; err != nil {
		return OutAuthDtos{}, fiber.ErrUnauthorized
	}

	userDtos := model.User{
		Id: tokenFound.UserId,
	}
	userFound := model.User{}
	if err := db.Preload("Role").Model(&userDtos).Where("id = ?", userDtos.Id).First(&userFound).Error; err != nil {
		return OutAuthDtos{}, err
	}

	outDtos := OutAuthDtos{
		Id:         userFound.Id,
		Username:   userFound.Username,
		Fullname:   userFound.Fullname,
		Email:      userFound.Email,
		Role:       userFound.Role.Name,
		VerifiedAt: userFound.VerifiedAt.String(),
		CreatedAt:  userFound.CreatedAt.String(),
	}

	return outDtos, nil
}

type OutAuthDtos struct {
	Id         int64  `json:"id"`
	Username   string `json:"username"`
	Fullname   string `json:"fullname"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	VerifiedAt string `json:"verified_at"`
	CreatedAt  string `json:"created_at"`
}
