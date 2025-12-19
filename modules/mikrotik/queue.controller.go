package mikrotik

import (
	"net/http"

	"mikrotik-api/config"

	"github.com/gin-gonic/gin"
)

func CreateSimpleQueue(c *gin.Context) {
	routerID := c.Param("id")
	userID := c.GetUint("user_id")

	var router Router
	if err := config.DB.First(&router, routerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}

	var body struct {
		Name   string `json:"name"`
		Target string `json:"target"`
		MaxLim string `json:"max_limit"`
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
		"/queue/simple/add",
		"=name="+body.Name,
		"=target="+body.Target,
		"=max-limit="+body.MaxLim,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	LogAction(userID, router.ID, "CREATE_QUEUE", body.Name)

	c.JSON(http.StatusCreated, gin.H{"message": "Queue created"})
}
