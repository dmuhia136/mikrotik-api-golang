package main

import (
	"mikrotik-api/config"
	"mikrotik-api/routes"

	"mikrotik-api/modules/auth"
	"mikrotik-api/modules/mikrotik"

	"github.com/gin-gonic/gin"
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
	routes.RegisterRoutes(r)

	println("🚀 MikroTik API running on port " + config.GetEnv("APP_PORT"))
	r.Run(":" + config.GetEnv("APP_PORT"))
}

