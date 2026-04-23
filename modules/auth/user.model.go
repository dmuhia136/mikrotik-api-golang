package auth

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(100)"`
	Email    string `json:"email" gorm:"type:varchar(255);uniqueIndex"`
	Password string `json:"-"`
	Role     string `json:"role" gorm:"type:varchar(50)"`
	Active   bool   `json:"active" gorm:"default:true"`
}
