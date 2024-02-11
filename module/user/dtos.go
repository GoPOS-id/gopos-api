package user

import "time"

type inUserDto struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	RoleId   uint   `json:"role_id"`
}

type outUserDto struct {
	Id         int64      `json:"id"`
	Username   string     `json:"username"`
	Fullname   string     `json:"fullname"`
	Email      string     `json:"email"`
	Role       string     `json:"role"`
	VerifiedAt *time.Time `json:"verified_at"`
	CreatedAt  time.Time  `json:"created_at"`
}

type outPaginateDto struct {
	Pagination interface{} `json:"pagination"`
	Users      interface{} `json:"users"`
}
