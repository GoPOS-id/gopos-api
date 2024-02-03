package auth

import (
	"github.com/GoPOS-id/gopos-api/api/model"
	"github.com/GoPOS-id/gopos-api/database"
	"github.com/GoPOS-id/gopos-api/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// login performs user authentication and generates a JWT token upon successful login.
// It retrieves the user details from the database based on the provided username,
// compares the hashed password with the input password, creates a JWT token,
// stores a session record in the database, and returns the authentication data.
// If any step fails, it returns an empty outAuthDtos and an error.
func (dtos *inAuthDtos) login() (outAuthDtos, error) {
	db := database.DbContext() // Get the database context

	// Find the user in the database by the provided username
	userDtos := model.User{
		Username: dtos.Username,
	}
	userFound := model.User{}
	if err := db.Model(&userDtos).Where("username = ?", userDtos.Username).First(&userFound).Error; err != nil {
		return outAuthDtos{}, err // Return an error if the user is not found in the database
	}

	// Compare the hashed password with the input password
	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(dtos.Password)); err != nil {
		return outAuthDtos{}, gorm.ErrRecordNotFound // Return an error if the password does not match
	}

	// Create a JWT token
	token, err := utils.CreateToken(userFound.Id, userFound.Username)
	if err != nil {
		return outAuthDtos{}, err // Return an error if token creation fails
	}

	// Map the session data and store it in the database
	session := mapSessionDto(userFound.Id, token)
	if err := db.Create(&session).Error; err != nil {
		return outAuthDtos{}, err // Return an error if session creation fails
	}

	outDtos := mapOutDto(userFound.Id, userFound.Username, token) // Map the output data for the successful login
	return outDtos, nil
}

// logout performs the user logout process by deleting the session record associated with the provided token.
// It retrieves the session details from the database based on the provided token,
// deletes the session record, and returns an error if any step fails.
func (dtos *outAuthDtos) logout() error {
	db := database.DbContext() // Get the database context

	// Find the session in the database by the provided token
	tokenDtos := model.Session{
		Token: dtos.Token,
	}
	tokenFound := model.Session{}
	if err := db.Model(&tokenDtos).Where("token = ?", tokenDtos.Token).First(&tokenFound).Error; err != nil {
		return err // Return an error if the session is not found in the database
	}

	// Delete the session record from the database
	if err := db.Delete(&tokenFound).Error; err != nil {
		return err // Return an error if session deletion fails
	}

	return nil // Return nil if the logout process is successful
}
