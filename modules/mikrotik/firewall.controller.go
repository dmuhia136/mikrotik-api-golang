package mikrotik

import (
	"net/http"

	"mikrotik-api/config"

	"github.com/gin-gonic/gin"
)

func AddFirewallRule(c *gin.Context) {
	routerID := c.Param("id")
	userID := c.GetUint("user_id")

	var router Router
	if err := config.DB.First(&router, routerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}

	var body struct {
		Chain    string `json:"chain"`
		Protocol string `json:"protocol"`
		Port     string `json:"port"`
		Action   string `json:"action"`
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
		"/ip/firewall/filter/add",
		"=chain="+body.Chain,
		"=protocol="+body.Protocol,
		"=dst-port="+body.Port,
		"=action="+body.Action,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	LogAction(userID, router.ID, "ADD_FIREWALL_RULE", body.Protocol+":"+body.Port)

	c.JSON(http.StatusCreated, gin.H{"message": "firewall rule added"})
}
