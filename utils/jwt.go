package utils

import (
	"strconv"
	"time"

	"github.com/GoPOS-id/gopos-api/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	Id       string
	Username string
	jwt.RegisteredClaims
}

var expires = time.Now().Add(config.JWT_EXP * time.Hour)

func CreateToken(dapet int64, Username string) (string, error) {
	Id := strconv.Itoa(int(dapet))
	claims := JWTClaims{
		Id,
		Username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "GoPOS",
		},
	}

	createToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := createToken.SignedString([]byte(config.JWT_SECRET_KEY))

	return token, err
}
