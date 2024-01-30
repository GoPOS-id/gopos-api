package auth

type inAuthDtos struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string
}

type outAuthDtos struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
