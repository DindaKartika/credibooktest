package models

type User struct {
	BaseModel
	Username string `json:"username" sql:"NOTNULL"`
	Password string `json:"password" sql:"NOTNULL"`
	IsAdmin  bool   `json:"is_admin"`
}

type UserResponse struct {
	BaseModel
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
}
