package user

import (
	"time"

	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/constant"
	"github.com/GoPOS-id/gopos-api/database"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// ! CREATE USER SERVICES
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
