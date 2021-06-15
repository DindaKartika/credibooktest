package models

type Transaction struct {
	BaseModel
	UserID int    `json:"user_id" sql:"NOTNULL"`
	Amount int    `json:"amount" sql:"NOTNULL"`
	Notes  string `json:"notes"`
	Type   string `json:"type"`
}
