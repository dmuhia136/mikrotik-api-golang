package auth

import "github.com/gin-gonic/gin"

func AuthRoutes(r *gin.RouterGroup) {
	r.POST("/register", Register)
	r.POST("/login", Login)
	r.POST("/logout", Logout)

}
