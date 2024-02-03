package user

import "time"

// inUserDto represents the input data structure for creating or updating a user.
// It includes fields for the user's ID, username, password, full name, email, and role ID.
type inUserDto struct {
	Id       int64  `json:"id"`       // User ID (used in update operations)
	Username string `json:"username"` // Username for the user
	Password string `json:"password"` // Password for the user
	Fullname string `json:"fullname"` // Full name of the user
	Email    string `json:"email"`    // Email address of the user
	RoleId   uint   `json:"role_id"`  // Role ID assigned to the user
}

// outUserDto represents the output data structure for retrieving user information.
// It includes fields for the user's ID, username, full name, email, role, verification timestamp, and creation timestamp.
type outUserDto struct {
	Id         int64      `json:"id"`         // User ID
	Username   string     `json:"username"`   // Username of the user
	Fullname   string     `json:"fullname"`   // Full name of the user
	Email      string     `json:"email"`      // Email address of the user
	Role       string     `json:"role"`       // Role assigned to the user
	VerifiedAt *time.Time `json:"verfied_at"` // Timestamp indicating when the user was verified
	CreatedAt  time.Time  `json:"created_at"` // Timestamp indicating when the user was created
}

// "current_page": page,       // Current page number
// "total_pages":  totalPages, // Total number of pages
// "total_data":   totalItems, // Total number of items
// "previous":     previous,   // Previous page number
// "next":         next,

type outPaginateDto struct {
	Pagination interface{} `json:"pagination"`
	Users      interface{} `json:"users"`
}
