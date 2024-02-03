package user

import (
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
// getProfileHandler handles the retrieval of user profile data and sends a successful data response.
// It retrieves the user profile data from the context locals, assuming it has been set by authentication middleware.
// Returns a 200 response with a success message and the user profile data.
func getProfileHandler(c *fiber.Ctx) error {
	locals := c.Locals("user").(middleware.OutAuthDtos)                       // Retrieve user profile data from the context locals
	return utils.SendDataResponse(c, "Success get profile data", locals, 200) // Send a successful data response with a 200 status code, success message, and user profile data
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
// @example {"pagination": {"current_page": 1, "total_pages": 5}, "users": [{"id": 1, "username": "user1"}, {"id": 2, "username": "user2"}]}
// getAllUsersHandler handles the retrieval of all users with pagination.
// It retrieves the page number from the query parameters, fetches users from the database
// with pagination, and sends a successful data response with user information and pagination details.
// Returns a 200 response with a success message, user data, and pagination information on success.
func getAllUsersHandler(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1) // Retrieve the page number from the query parameters (default to 1 if not specified)
	offset := (page - 1) * 25     // Calculate the offset based on the page number and limit (25 users per page)

	db := database.DbContext() // Get the database context

	// Fetch users from the database with pagination and preload the "Role" relationship
	var usersDtos []model.User
	if err := db.Model(&model.User{}).Preload("Role").Offset(offset).Limit(25).Find(&usersDtos).Error; err != nil {
		return utils.SendResponse(c, err.Error(), 500) // Handle database error and send a 500 response
	}

	// Perform pagination on the retrieved users
	outUserDtos, previous, next, totalItems, totalPages := handleUsersPagination(db, usersDtos, page)

	// Create pagination map
	paginate := mapPaginationUsers(page, previous, next, totalItems, totalPages)

	// Combine user data and pagination information into a response map
	outDto := mapOutAllUserDto(outUserDtos, paginate)

	// Send a successful data response with a 200 status code, success message, user data, and pagination information
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
// createUserHandler handles the creation of a new user.
// It parses the request body to extract user data, performs validation,
// checks permission to add a user with a specific role, creates the user,
// and sends an appropriate response based on the result.
// Returns a 200 response with a success message and the created user on success.
func createUserHandler(c *fiber.Ctx) error {
	dtos := new(inUserDto)                              // Parse the request body to extract user data
	locals := c.Locals("user").(middleware.OutAuthDtos) // Retrieve user role from the context locals

	// Parse the request body
	if err := c.BodyParser(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 409)
	}

	// Validate user data
	if err := handleCreateValidator(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 400)
	}

	// Check permission to add a user with a specific role
	if !handleAddRolePermission(dtos.RoleId, locals.Role) {
		return utils.SendResponse(c, "Can't create role higher than your currnet role", 400)
	}

	// Create the user
	user, err := dtos.Create()
	if err != nil {
		return handleErrCreateResponse(c, err) // Handle errors during user creation and send an appropriate response
	}

	return utils.SendDataResponse(c, "Create User Successfully", user, 200) // Send a successful data response with a 200 status code, success message, and the created user
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
// updateUserHandler handles the update of user data.
// It parses the request body to extract updated user data, performs validation,
// checks permission to update the user with a specific role, updates the user,
// and sends an appropriate response based on the result.
// Returns a 200 response with a success message and the updated user data on success.
func updateUserHandler(c *fiber.Ctx) error {
	dtos := new(inUserDto)                              // Parse the request body to extract updated user data
	locals := c.Locals("user").(middleware.OutAuthDtos) // Retrieve user role from the context locals

	// Parse the request body
	if err := c.BodyParser(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 409) // Handle parsing error and send a 409 response
	}

	// Validate updated user data
	if err := handleUpdateUserValidator(dtos); err != nil {
		return utils.SendResponse(c, err.Error(), 409) // Handle validation error and send a 409 response
	}

	// Check permission to update the user with a specific role
	if !handleAddRolePermission(dtos.RoleId, locals.Role) {
		return utils.SendResponse(c, "Can't change role higher than your currnet role", 400) // Send a 400 response if the user attempts to change to a higher role than their current role
	}

	// Update the user
	user, errDtos := dtos.Update()
	if errDtos != nil {
		return handleErrUpdateUser(c, errDtos) // Handle errors during user update and send an appropriate response
	}

	// Send a successful data response with a 200 status code, success message, and the updated user
	return utils.SendDataResponse(c, "Success update data users", user, 200)
}
