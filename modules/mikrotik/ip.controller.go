package mikrotik

import (
	"net/http"

	"mikrotik-api/config"

	"github.com/gin-gonic/gin"
)

func AddIPAddress(c *gin.Context) {
	routerID := c.Param("id")
	userID := c.GetUint("user_id")

	var router Router
	if err := config.DB.First(&router, routerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}

	var body struct {
		Address   string `json:"address"`
		Interface string `json:"interface"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := Connect(router)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "connection failed"})
		return
	}
	defer client.Close()

	_, err = client.Run(
		"/ip/address/add",
		"=address="+body.Address,
		"=interface="+body.Interface,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	LogAction(userID, router.ID, "ADD_IP_ADDRESS", body.Address)

	c.JSON(http.StatusCreated, gin.H{"message": "IP address added"})
}
