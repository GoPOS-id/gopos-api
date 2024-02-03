package user

import (
	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/constant"
	"github.com/GoPOS-id/gopos-api/utils"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// handleCreateValidator performs validation on the input data for creating a new user.
// It checks if the provided fields (username, password, fullname, email, and role ID)
// meet the required length and presence criteria. Returns an error if validation fails or nil otherwise.
func handleCreateValidator(dtos *inUserDto) error {
	return validation.ValidateStruct(dtos,
		validation.Field(&dtos.Username, validation.Required, validation.Length(8, 25)),
		validation.Field(&dtos.Password, validation.Required, validation.Length(8, 25)),
		validation.Field(&dtos.Fullname, validation.Required, validation.Length(8, 30)),
		validation.Field(&dtos.Email, validation.Required),
		validation.Field(&dtos.RoleId, validation.Required),
	)
}

// handleAddRolePermission checks if a given role ID and role string satisfy permission conditions
// for adding a new role. It returns true if the conditions are met, and false otherwise.
// The conditions are:
// - Role ID is 1 or 2 (indicating administrator or manager).
// - If the role ID is 1 or 2, the role string must be "Operator".
// Returns true for other role IDs.
func handleAddRolePermission(roleId uint, role string) bool {
	if roleId == 1 || roleId == 2 {
		if role != constant.Role_Operator {
			return false // If the role ID is 1 or 2, the role string must be "Operator" for permission.
		} else {
			return true // Permission is granted if the role ID is 1 or 2 and the role string is "Operator".
		}
	}
	return true // Permission is granted for other role IDs.
}

// handleErrCreateResponse handles specific errors that may occur during user creation
// and sends an appropriate HTTP response.
// If the error is a fiber.ErrBadRequest, it sends a 400 response with a specific error message.
// For other errors, it sends a generic 500 response with the error message.
func handleErrCreateResponse(c *fiber.Ctx, err error) error {
	switch err {
	case fiber.ErrBadRequest:
		return utils.SendResponse(c, "Username already exists", 400) // Send a 400 response with a specific error message for duplicate usernames.
	default:
		return utils.SendResponse(c, err.Error(), 500) // Send a generic 500 response with the error message for other errors.
	}
}

// handleUsersPagination performs pagination for a list of users and transforms them into outUserDto objects.
// It takes a GORM database connection, a slice of model.User, and the current page number as input.
// It calculates the total number of items, total pages, previous and next page numbers,
// and maps the users to outUserDto objects.
// Returns the mapped outUserDto objects, previous and next page numbers, total number of items, and total pages.
func handleUsersPagination(db *gorm.DB, dtos []model.User, page int) ([]outUserDto, int, int, int64, int64) {
	// Calculate the total number of items in the database table
	var totalItems int64
	db.Model(&model.User{}).Count(&totalItems)

	// Calculate the total number of pages based on the items per page (25)
	totalPages := (totalItems + int64(25) - 1) / int64(25)

	// Calculate the previous and next page numbers
	var previous int = page - 1
	var next int = page + 1

	if previous <= 1 {
		previous = 1
	}

	if next >= int(totalPages) {
		next = int(totalPages)
	}

	// Map the users to outUserDto objects
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

	return outUserDtos, previous, next, totalItems, totalPages // Return the mapped outUserDto objects, previous and next page numbers, total number of items, and total pages
}

// mapPaginationUsers creates a map representing pagination information for users.
// It takes the current page number, previous and next page numbers, total number of items,
// and total number of pages as input and returns a fiber.Map containing this pagination information.
func mapPaginationUsers(page int, previous int, next int, totalItems int64, totalPages int64) interface{} {
	return fiber.Map{
		"current_page": page,       // Current page number
		"total_pages":  totalPages, // Total number of pages
		"total_data":   totalItems, // Total number of items
		"previous":     previous,   // Previous page number
		"next":         next,       // Next page number
	}
}

// mapOutAllUserDto creates a map combining a list of outUserDto objects and pagination information.
// It takes a slice of outUserDto representing user data and pagination information as input,
// and returns a fiber.Map containing the combined data.
func mapOutAllUserDto(userDto []outUserDto, paginate interface{}) interface{} {
	// Create and return a fiber.Map with user data and pagination information
	return outPaginateDto{
		Users:      userDto,
		Pagination: paginate,
	}
}

// handleUpdateUserValidator performs validation on the input data for updating a user.
// It checks if the provided fields (ID, username, password, fullname, email, and role ID)
// meet the required length and presence criteria. Returns an error if validation fails or nil otherwise.
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

// handleErrUpdateUser handles specific errors that may occur during the user update process
// and sends an appropriate HTTP response.
// If the error is gorm.ErrRecordNotFound, it sends a 404 response with a specific error message.
// For other errors, it sends a generic 500 response with the error message.
func handleErrUpdateUser(c *fiber.Ctx, err error) error {
	switch err {
	case gorm.ErrRecordNotFound:
		return utils.SendResponse(c, "User not found", 404)
	default:
		return utils.SendResponse(c, err.Error(), 500)
	}
}
