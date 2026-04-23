package routes

import (
	"mikrotik-api/modules/auth"
	"mikrotik-api/modules/middleware"

	"github.com/gin-gonic/gin"
	"mikrotik-api/modules/mikrotik"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Public
	auth.AuthRoutes(api)

	// Protected
	protected := api.Group("/")
	protected.Use(middleware.JWTAuth())

	mikrotik.MikroTikRoutes(api)

	// Example protected test route
	protected.GET("/me", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"user_id": c.GetUint("user_id"),
			"role":    c.GetString("role"),
		})
	})

	// Admin only
	admin := protected.Group("/admin")
	admin.Use(middleware.RequireRole("admin"))
	admin.GET("/dashboard", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Admin access granted"})
	})

	admin.POST("/users", auth.AdminCreateUser)
	admin.GET("/users", auth.AdminListUsers)
	admin.PUT("/users/:id/role", auth.AdminUpdateUserRole)
	admin.PUT("/users/:id/status", auth.AdminSetUserStatus)
}
