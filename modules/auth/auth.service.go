package auth

import (
	"errors"
	"mikrotik-api/config"
	"mikrotik-api/utils"
)

func RegisterUser(name, email, password string) (*User, error) {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := User{
		Name:     name,
		Email:    email,
		Password: hash,
		Role:     "admin", // first user = admin
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func LoginUser(email, password string) (*User, error) {
	var user User

	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return &user, nil
}
