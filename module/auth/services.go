package auth

import (
	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/database"
	"github.com/GoPOS-id/gopos-api/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (dtos *inAuthDtos) login() (outAuthDtos, error) {
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

	session := mapSessionDto(userFound.Id, token)
	if err := db.Create(&session).Error; err != nil {
		return outAuthDtos{}, err
	}

	outDtos := mapOutDto(userFound.Id, userFound.Username, token)
	return outDtos, nil
}

func (dtos *outAuthDtos) logout() error {
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
