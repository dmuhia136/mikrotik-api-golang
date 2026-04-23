package main

import (
	"time"

	"mikrotik-api/config"
	"mikrotik-api/routes"

	"mikrotik-api/modules/auth"
	"mikrotik-api/modules/mikrotik"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	if err := config.DB.AutoMigrate(
		&auth.User{},
		&mikrotik.Router{},
		&mikrotik.AuditLog{},
	); err != nil {
		panic("Migration failed: " + err.Error())
	}

	r := gin.Default()

	// 🔥 CORS CONFIG (THIS FIXES YOUR ISSUE)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.RegisterRoutes(r)

	r.Run(":" + config.GetEnv("APP_PORT"))
}
