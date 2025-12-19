package mikrotik

import (
	"net/http"

	"mikrotik-api/config"

	"github.com/gin-gonic/gin"
)

func GetInterfaces(c *gin.Context) {
	routerID := c.Param("id")

	var router Router
	if err := config.DB.First(&router, routerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}

	client, err := Connect(router)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "connection failed"})
		return
	}
	defer client.Close()

	reply, err := client.Run("/interface/print")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reply.Re)
}


func CreateRouter(c *gin.Context) {
	var router Router

	if err := c.ShouldBindJSON(&router); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if router.Port == 0 {
		router.Port = 8729
		router.UseTLS = true
	}

	if err := config.DB.Create(&router).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, router)
}

func ListRouters(c *gin.Context) {
	var routers []Router
	config.DB.Find(&routers)
	c.JSON(http.StatusOK, routers)
}