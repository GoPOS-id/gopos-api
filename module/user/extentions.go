package user

import (
	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/constant"
	"github.com/GoPOS-id/gopos-api/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func handleCreateValidator(dtos *inUserDto) error {
	return validation.ValidateStruct(dtos,
		validation.Field(&dtos.Username, validation.Required, validation.Length(8, 25)),
		validation.Field(&dtos.Password, validation.Required, validation.Length(8, 25)),
		validation.Field(&dtos.Fullname, validation.Required, validation.Length(8, 30)),
		validation.Field(&dtos.Email, validation.Required),
		validation.Field(&dtos.RoleId, validation.Required),
	)
}

func handleAddRolePermission(roleId uint, role string) bool {
	if roleId == 1 || roleId == 2 {
		if role != constant.Role_Operator {
			return false
		} else {
			return true
		}
	}
	return true
}

func handleUsersPagination(db *gorm.DB, category string, dtos []model.User, page int, limit int) ([]outUserDto, int, int, int64, int64) {
	var totalItems int64
	switch category {
	case "operator":
		db.Model(&model.User{}).Where("role_id = ?", 1).Count(&totalItems)
	case "administrator":
		db.Model(&model.User{}).Where("role_id = ?", 2).Count(&totalItems)
	case "cashier":
		db.Model(&model.User{}).Where("role_id = ?", 3).Count(&totalItems)
	default:
		db.Model(&model.User{}).Count(&totalItems)
	}

	totalPages := (totalItems + int64(limit) - 1) / int64(limit)

	var previous int = page - 1
	var next int = page + 1

	if previous <= 1 {
		previous = 1
	}

	if next >= int(totalPages) {
		next = int(totalPages)
	}

	var outUserDtos []outUserDto
	for _, u := range dtos {
		userDto := outUserDto{
			Id:         u.Id,
			Username:   u.Username,
			Fullname:   u.Fullname,
			Email:      u.Email,
			Role:       u.Role.Name,
			VerifiedAt: u.VerifiedAt,
			CreatedAt:  u.CreatedAt,
		}
		outUserDtos = append(outUserDtos, userDto)
	}

	return outUserDtos, previous, next, totalItems, totalPages
}

func mapPaginationUsers(page int, previous int, next int, totalItems int64, totalPages int64) interface{} {
	return fiber.Map{
		"current_page": page,
		"total_pages":  totalPages,
		"total_data":   totalItems,
		"previous":     previous,
		"next":         next,
	}
}

func mapOutAllUserDto(userDto []outUserDto, paginate interface{}) interface{} {
	return outPaginateDto{
		Users:      userDto,
		Pagination: paginate,
	}
}

func handleUpdateUserValidator(dtos *inUserDto) error {
	return validation.ValidateStruct(dtos,
		validation.Field(&dtos.Id, validation.Required),
		validation.Field(&dtos.Username, validation.Required, validation.Length(8, 25)),
		validation.Field(&dtos.Password, validation.Length(8, 25)),
		validation.Field(&dtos.Fullname, validation.Required, validation.Length(8, 30)),
		validation.Field(&dtos.Email, validation.Required),
		validation.Field(&dtos.RoleId, validation.Required),
	)
}

func handleErr(c *fiber.Ctx, err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return utils.SendResponse(c, "User not found", 404)
	case fiber.ErrBadRequest:
		return utils.SendResponse(c, "Username already exists", 400)
	case fiber.ErrNotAcceptable:
		return utils.SendResponse(c, "Oops!, You dont have permission to change the roles!", 400)
	default:
		return utils.SendResponse(c, err.Error(), 500)
	}
}

func handleDeleteUserValidator(dtos *inUserDto) error {
	return validation.ValidateStruct(dtos,
		validation.Field(&dtos.Id, validation.Required),
	)
}
