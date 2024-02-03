package auth

// inAuthDtos represents the input data structure for authentication-related operations.
// It includes fields for the username, password, and an optional token.
type inAuthDtos struct {
	Username string `json:"username"` // Username for authentication
	Password string `json:"password"` // Password for authentication
}

// outAuthDtos represents the output data structure for authentication-related operations.
// It includes fields for the user ID, username, and a token.
type outAuthDtos struct {
	Id       int64  `json:"id"`       // User ID
	Username string `json:"username"` // Username of the authenticated user
	Token    string `json:"token"`    // JWT token for authentication
}
