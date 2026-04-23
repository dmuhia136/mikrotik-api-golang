package auth

import (
	"errors"
	"mikrotik-api/config"
	"mikrotik-api/utils"

	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"

)

func CreateUser(name, email, password, role string) (*User, error) {
	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := User{
		Name:     name,
		Email:    email,
		Password: hash,
		Role:     role,
		Active:   true,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func ListUsers() ([]User, error) {
	var users []User
	err := config.DB.Find(&users).Error
	return users, err
}

func UpdateUserRole(userID uint, role string) error {
	return config.DB.Model(&User{}).
		Where("id = ?", userID).
		Update("role", role).Error
}

func SetUserStatus(userID uint, active bool) error {
	return config.DB.Model(&User{}).
		Where("id = ?", userID).
		Update("active", active).Error
}

func GetActiveUserByID(id uint) (*User, error) {
	var user User
	err := config.DB.Where("id = ? AND active = ?", id, true).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found or inactive")
	}
	return &user, nil
}



func AdminCreateUser(c *gin.Context) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := CreateUser(body.Name, body.Email, body.Password, body.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func AdminListUsers(c *gin.Context) {
	users, err := ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func AdminUpdateUserRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var body struct {
		Role string `json:"role"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := UpdateUserRole(uint(id), body.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role updated"})
}

func AdminSetUserStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var body struct {
		Active bool `json:"active"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := SetUserStatus(uint(id), body.Active); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user status updated"})
}