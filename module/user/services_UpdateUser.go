package user

import (
	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/database"
	"golang.org/x/crypto/bcrypt"
)

// ! UPDATE USER SERVICES
func (dto *inUserDto) Update() (outUserDto, error) {
	db := database.DbContext()
	userDtos := model.User{
		Id: dto.Id,
	}
	userFound := model.User{}
	if err := db.Preload("Role").Model(&userDtos).Where("id = ?", userDtos.Id).First(&userFound).Error; err != nil {
		return outUserDto{}, err
	}
	if dto.Password != "" {
		passwd, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
		userFound.Password = string(passwd)
	}
	userFound.Fullname = dto.Fullname
	userFound.Email = dto.Email
	userFound.RoleId = dto.RoleId

	if err := db.Save(&userFound).Error; err != nil {
		return outUserDto{}, err
	}

	outDto := outUserDto{
		Id:         userFound.Id,
		Username:   userFound.Username,
		Fullname:   userFound.Fullname,
		Email:      userFound.Email,
		Role:       userFound.Role.Name,
		VerifiedAt: userFound.VerifiedAt,
		CreatedAt:  userFound.CreatedAt,
	}

	return outDto, nil
}
