package auth

import (
	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/database"
)

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
