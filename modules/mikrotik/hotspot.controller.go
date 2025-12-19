package mikrotik

import (
	"net/http"

	"mikrotik-api/config"

	"github.com/gin-gonic/gin"
)

func CreateHotspotUser(c *gin.Context) {
	routerID := c.Param("id")
	userID := c.GetUint("user_id")

	var router Router
	if err := config.DB.First(&router, routerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}

	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Profile  string `json:"profile"`
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
		"/ip/hotspot/user/add",
		"=name="+body.Username,
		"=password="+body.Password,
		"=profile="+body.Profile,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	LogAction(userID, router.ID, "CREATE_HOTSPOT_USER", body.Username)

	c.JSON(http.StatusCreated, gin.H{"message": "Hotspot user created"})
}
