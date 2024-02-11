package user

import (
	"strconv"

	"github.com/GoPOS-id/gopos-api/api/middleware"
	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/database"
	"github.com/GoPOS-id/gopos-api/utils"
	"github.com/gofiber/fiber/v2"
)

// @Tags User
// @Summary Get Profile
// @Description Retrieves the profile data for the authenticated user.
// @ID getProfile
// @Produce json
// @Param Authorization header string true "Bearer <access_token>" default("Bearer ")
// @Success 200 {object} utils.DataResponse{data=[]outUserDto} "User profile data"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Router /user [get]
func getProfileHandler(c *fiber.Ctx) error {
	locals := c.Locals("user").(middleware.OutAuthDtos)
	return utils.SendDataResponse(c, "Success get profile data", locals, 200)
}

// @Tags User
// @Summary Get All Users
// @Description Retrieves a list of users with pagination support.
// @ID getAllUsers
// @Produce json
// @Param Authorization header string true "Bearer <access_token>" default("Bearer ")
// @Param page query int false "Page number for pagination (default: 1)"
// @Success 200 {object} utils.DataResponse{data=[]outPaginateDto} "List of users with pagination"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 500 {object} utils.Response "Internal Server Error"
// @Router /user/all [get]
func getAllUsersHandler(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := 20
	offset := (page - 1) * limit
	category := c.Query("category", "all")

	db := database.DbContext()

	var usersDtos []model.User
	switch category {
	case "operator":
		if err := db.Model(&model.User{}).Preload("Role").Where("role_id = ?", 1).Offset(offset).Limit(limit).Order("role_id asc").Find(&usersDtos).Error; err != nil {
			return utils.SendResponse(c, err.Error(), 500)
		}
	case "administrator":
		if err := db.Model(&model.User{}).Preload("Role").Where("role_id = ?", 2).Offset(offset).Limit(limit).Order("role_id asc").Find(&usersDtos).Error; err != nil {
			return utils.SendResponse(c, err.Error(), 500)
		}
	case "cashier":
		if err := db.Model(&model.User{}).Preload("Role").Where("role_id = ?", 3).Offset(offset).Limit(limit).Order("role_id asc").Find(&usersDtos).Error; err != nil {
			return utils.SendResponse(c, err.Error(), 500)
		}
	default:
		if err := db.Model(&model.User{}).Preload("Role").Offset(offset).Limit(limit).Order("role_id asc").Find(&usersDtos).Error; err != nil {
			return utils.SendResponse(c, err.Error(), 500)
		}
	}

	outUserDtos, previous, next, totalItems, totalPages := handleUsersPagination(db, category, usersDtos, page, limit)

	paginate := mapPaginationUsers(page, previous, next, totalItems, totalPages)

	outDto := mapOutAllUserDto(outUserDtos, paginate)

	return utils.SendDataResponse(c, "Success show all data users", outDto, 200)
}

// @Summary Create a new user
// @Description Create a new user with the provided data
// @ID createUser
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>" default("Bearer ")
// @Param user body inUserDto true "User data"
// @Success 200 {object} utils.DataResponse{data=inUserDto} "User created successfully"
// @Failure 400 {object} utils.Response "Bad Request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 409 {object} utils.Response "Conflict"
// @Failure 500 {object} utils.Response "Internal Server Error"
// @Router /user [post]
func createUserHandler(c *fiber.Ctx) error {
	dtos := new(inUserDto)
	locals := c.Locals("user").(middleware.OutAuthDtos)

	if err := c.BodyParser(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 409)
	}

	if err := handleCreateValidator(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 400)
	}

	if !handleAddRolePermission(dtos.RoleId, locals.Role) {
		return utils.SendResponse(c, "Can't create role higher than your currnet role", 400)
	}

	user, err := dtos.Create()
	if err != nil {
		return handleErr(c, err)
	}

	return utils.SendDataResponse(c, "Create User Successfully", user, 200)
}

// @Summary Update user data
// @Description Update user data with the provided information
// @ID updateUser
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>" default("Bearer ")
// @Param user body inUserDto true "Updated user data"
// @Success 200 {object} utils.DataResponse{data=inUserDto} "User updated successfully"
// @Failure 400 {object} utils.Response "Bad Request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 409 {object} utils.Response "Conflict"
// @Failure 500 {object} utils.Response "Internal Server Error"
// @Router /user [patch]
func updateUserHandler(c *fiber.Ctx) error {
	dtos := new(inUserDto)
	locals := c.Locals("user").(middleware.OutAuthDtos)

	if err := c.BodyParser(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 409)
	}

	if err := handleUpdateUserValidator(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 400)
	}

	if !handleAddRolePermission(dtos.RoleId, locals.Role) {
		return utils.SendResponse(c, "Can't change role higher than your currnet role", 400)
	}

	user, errDtos := dtos.Update(locals)
	if errDtos != nil {
		return handleErr(c, errDtos)
	}

	return utils.SendDataResponse(c, "Success update data users", user, 200)
}

// @Summary Delete user data
// @Description Delete user data with the provided information
// @ID deleteUser
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>" default("Bearer ")
// @Param user body inUserDto true "Delete user data"
// @Success 200 {object} utils.Response "User updated successfully"
// @Failure 400 {object} utils.Response "Bad Request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 409 {object} utils.Response "Conflict"
// @Failure 500 {object} utils.Response "Internal Server Error"
// @Router /user [delete]
func deleteUserHandler(c *fiber.Ctx) error {
	dtos := new(inUserDto)
	locals := c.Locals("user").(middleware.OutAuthDtos)

	if err := c.BodyParser(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 409)
	}

	if err := handleDeleteUserValidator(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 400)
	}

	if !handleAddRolePermission(dtos.RoleId, locals.Role) {
		return utils.SendResponse(c, "You cant delete user with role higher than you", 400)
	}

	if err := dtos.Delete(); err != nil {
		return handleErr(c, err)
	}
	return utils.SendResponse(c, "Delete user successfully", 200)
}

// @Summary Show User by ID
// @Description Show user data by ID with the provided information
// @ID showIdUser
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer <access_token>" default("Bearer ")
// @Success 200 {object} utils.Response "Show User updated successfully"
// @Failure 400 {object} utils.Response "Bad Request"
// @Failure 401 {object} utils.Response "Unauthorized"
// @Failure 404 {object} utils.Response "User not found"
// @Failure 409 {object} utils.Response "Conflict"
// @Failure 500 {object} utils.Response "Internal Server Error"
// @Router /user/view/{id} [get]
func viewUserHandler(c *fiber.Ctx) error {
	params := c.Params("userid")
	userid, err := strconv.Atoi(params)
	if err != nil {
		return utils.SendResponse(c, "User ID not found", 404)
	}
	dtos := inUserDto{Id: int64(userid)}

	user, errors := dtos.GetProfile()
	if errors != nil {
		return handleErr(c, errors)
	}

	return utils.SendDataResponse(c, "Show ID "+params+" successfully", user, 200)

}
