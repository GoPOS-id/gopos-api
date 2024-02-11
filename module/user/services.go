package user

import (
	"time"

	"github.com/GoPOS-id/gopos-api/api/middleware"
	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/constant"
	"github.com/GoPOS-id/gopos-api/database"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (dto *inUserDto) Create() (outUserDto, error) {
	db := database.DbContext()

	userDto := model.User{
		Username: dto.Username,
	}
	foundDto := model.User{}
	if err := db.Model(&userDto).Where("username = ?", userDto.Username).First(&foundDto).Error; err != gorm.ErrRecordNotFound {
		return outUserDto{}, fiber.ErrBadRequest
	}

	guid := xid.New()
	passwd, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
	currentTime := time.Now()
	pointerTime := &currentTime
	create := model.User{
		Id:         guid.Time().Unix(),
		Username:   dto.Username,
		Password:   string(passwd),
		Fullname:   dto.Fullname,
		Email:      dto.Email,
		RoleId:     dto.RoleId,
		VerifiedAt: pointerTime,
	}

	if err := db.Create(&create).Error; err != nil {
		return outUserDto{}, err
	}
	var Role string
	if dto.RoleId == 1 {
		Role = constant.Role_Operator
	} else if dto.RoleId == 2 {
		Role = constant.Role_Adminstrator
	} else if dto.RoleId == 3 {
		Role = constant.Role_Cashier
	}

	outDto := outUserDto{
		Id:         create.Id,
		Username:   create.Username,
		Fullname:   create.Fullname,
		Email:      create.Email,
		Role:       Role,
		VerifiedAt: pointerTime,
		CreatedAt:  currentTime,
	}

	return outDto, nil
}

func (dto *inUserDto) Update(locals middleware.OutAuthDtos) (outUserDto, error) {
	db := database.DbContext()
	userDtos := model.User{
		Id: dto.Id,
	}
	userFound := model.User{}
	if err := db.Preload("Role").Model(&userDtos).Where("id = ?", userDtos.Id).First(&userFound).Error; err != nil {
		return outUserDto{}, err
	}

	if locals.Role != "operator" {
		if dto.RoleId != userFound.RoleId {
			return outUserDto{}, fiber.ErrNotAcceptable
		}
	}

	var passwd string = userFound.Password
	if dto.Password != "" {
		genereate, _ := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.DefaultCost)
		passwd = string(genereate)
	}

	if err := db.Model(&userFound).Updates(
		model.User{
			Fullname: dto.Fullname,
			Email:    dto.Email,
			Password: passwd,
			RoleId:   dto.RoleId,
		}).Error; err != nil {
		return outUserDto{}, err
	}
	roleDtos := model.Role{Id: dto.RoleId}
	foundDtos := model.Role{}
	if err := db.Model(&roleDtos).Where("id = ?", roleDtos.Id).First(&foundDtos).Error; err != nil {
		return outUserDto{}, err
	}

	outDto := outUserDto{
		Id:         userFound.Id,
		Username:   userFound.Username,
		Fullname:   userFound.Fullname,
		Email:      userFound.Email,
		Role:       foundDtos.Name,
		VerifiedAt: userFound.VerifiedAt,
		CreatedAt:  userFound.CreatedAt,
	}

	return outDto, nil
}

func (dto *inUserDto) Delete() error {
	db := database.DbContext()

	userDtos := model.User{
		Id: dto.Id,
	}
	userFound := model.User{}
	if err := db.Model(&userDtos).Where("id = ?", userDtos.Id).First(&userFound).Error; err != nil {
		return err
	}

	if err := db.Model(&userFound).Delete(&userFound).Error; err != nil {
		return err
	}

	return nil
}

func (dto *inUserDto) GetProfile() (outUserDto, error) {
	db := database.DbContext()

	userDto := model.User{Id: dto.Id}
	foundDto := model.User{}
	if err := db.Preload("Role").Model(&userDto).Where("id = ?", userDto.Id).First(&foundDto).Error; err != nil {
		return outUserDto{}, err
	}

	res := outUserDto{
		Id:         foundDto.Id,
		Username:   foundDto.Username,
		Fullname:   foundDto.Fullname,
		Email:      foundDto.Email,
		Role:       foundDto.Role.Name,
		VerifiedAt: foundDto.VerifiedAt,
		CreatedAt:  foundDto.CreatedAt,
	}

	return res, nil
}
