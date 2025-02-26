package models

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Password string `json:"-"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar_url"`
}

type UserDataRequest struct {
	Username string `gorm:"unique" json:"username"`
	Password string `json:"password"`
}
