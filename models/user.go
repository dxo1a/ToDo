package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"-"`
}

type UserDataRequest struct {
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
}
