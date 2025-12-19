package mikrotik

import (
	"net/http"
	"strconv"

	"mikrotik-api/config"

	"github.com/gin-gonic/gin"
)

// ---------- PPP PROFILE ----------
// Creates a PPP profile (used by PPPoE secrets)
func CreatePPPProfile(c *gin.Context) {
	routerID := c.Param("id")
	userID := c.GetUint("user_id")

	var router Router
	if err := config.DB.First(&router, routerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}

	var body struct {
		Name         string `json:"name"`          // e.g. "10Mbps"
		LocalAddress string `json:"local_address"` // optional
		RemotePool   string `json:"remote_pool"`   // optional
		RateLimit    string `json:"rate_limit"`    // e.g. "10M/10M" (RouterOS format) optional
		DNS1         string `json:"dns1"`          // optional
		DNS2         string `json:"dns2"`          // optional
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

	args := []string{
		"/ppp/profile/add",
		"=name=" + body.Name,
	}

	// Optional params
	if body.LocalAddress != "" {
		args = append(args, "=local-address="+body.LocalAddress)
	}
	if body.RemotePool != "" {
		args = append(args, "=remote-address="+body.RemotePool)
	}
	if body.RateLimit != "" {
		args = append(args, "=rate-limit="+body.RateLimit)
	}
	if body.DNS1 != "" {
		args = append(args, "=dns-server="+body.DNS1)
		if body.DNS2 != "" {
			args[len(args)-1] = args[len(args)-1] + "," + body.DNS2
		}
	}

	_, err = client.RunArgs(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	LogAction(userID, router.ID, "CREATE_PPP_PROFILE", body.Name)
	c.JSON(http.StatusCreated, gin.H{"message": "PPP profile created"})
}

func ListPPPProfiles(c *gin.Context) {
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

	reply, err := client.Run("/ppp/profile/print")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reply.Re)
}

// ---------- HOTSPOT PROFILE ----------
func CreateHotspotUserProfile(c *gin.Context) {
	routerID := c.Param("id")
	userID := c.GetUint("user_id")

	var router Router
	if err := config.DB.First(&router, routerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "router not found"})
		return
	}

	var body struct {
		Name      string `json:"name"`       // e.g. "Daily-1Mbps"
		RateLimit string `json:"rate_limit"` // e.g. "1M/1M"
		Shared    int    `json:"shared"`     // e.g. 1, 2, 3
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Shared <= 0 {
		body.Shared = 1
	}

	client, err := Connect(router)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "connection failed"})
		return
	}
	defer client.Close()

	args := []string{
		"/ip/hotspot/user/profile/add",
		"=name=" + body.Name,
		"=shared-users=" + strconv.Itoa(body.Shared),
	}
	if body.RateLimit != "" {
		args = append(args, "=rate-limit="+body.RateLimit)
	}

	_, err = client.RunArgs(args)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	LogAction(userID, router.ID, "CREATE_HOTSPOT_PROFILE", body.Name)
	c.JSON(http.StatusCreated, gin.H{"message": "Hotspot user profile created"})
}

func ListHotspotUserProfiles(c *gin.Context) {
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

	reply, err := client.Run("/ip/hotspot/user/profile/print")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reply.Re)
}
