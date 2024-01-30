package auth

import (
	"strings"

	"github.com/GoPOS-id/gopos-api/database"
	"github.com/GoPOS-id/gopos-api/model"
	"github.com/GoPOS-id/gopos-api/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ! LOGIN SERVICES
func (dtos *inAuthDtos) Login() (outAuthDtos, error) {
	db := database.DbContext()
	userDtos := model.User{
		Username: dtos.Username,
	}
	userFound := model.User{}
	if err := db.Model(&userDtos).Where("username = ?", userDtos.Username).First(&userFound).Error; err != nil {
		return outAuthDtos{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(dtos.Password)); err != nil {
		return outAuthDtos{}, gorm.ErrRecordNotFound
	}

	token, err := utils.CreateToken(userFound.Id, userFound.Username)

	if err != nil {
		return outAuthDtos{}, err
	}

	session := model.Session{
		UserId: userFound.Id,
		Token:  strings.Split(token, ".")[2],
	}

	if err := db.Create(&session).Error; err != nil {
		return outAuthDtos{}, err
	}

	outDtos := outAuthDtos{
		Id:       userFound.Id,
		Username: userFound.Username,
		Token:    token,
	}

	return outDtos, nil
}

// ! LOGOUT SERVICES
func (dtos *inAuthDtos) Logout() error {
	db := database.DbContext()
	tokenDtos := model.Session{
		Token: dtos.Token,
	}
	tokenFound := model.Session{}
	if err := db.Model(&tokenDtos).Where("token = ?", tokenDtos.Token).First(&tokenFound).Error; err != nil {
		return err
	}

	if err := db.Delete(&tokenFound).Error; err != nil {
		return err
	}

	return nil
}
