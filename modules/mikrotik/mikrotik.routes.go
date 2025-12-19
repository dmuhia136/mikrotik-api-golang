package mikrotik

import (
	"mikrotik-api/modules/middleware"

	"github.com/gin-gonic/gin"
)

func MikroTikRoutes(r *gin.RouterGroup) {
	r.Use(middleware.JWTAuth())

	// Router registry
	admin := r.Group("/routers")
	admin.Use(middleware.RequireRole("admin"))
	admin.POST("/", CreateRouter)
	admin.GET("/", ListRouters)

	// Device control (admin + operator)
	control := r.Group("/routers/:id")
	control.Use(middleware.RequireRole("admin", "operator"))

	control.POST("/firewall", AddFirewallRule)
	control.POST("/ip", AddIPAddress)

	control.POST("/pppoe", CreatePPPoEUser)
	control.POST("/hotspot", CreateHotspotUser)
	control.POST("/queue", CreateSimpleQueue)

	control.POST("/ppp-profiles", CreatePPPProfile)
	control.GET("/ppp-profiles", ListPPPProfiles)

	control.POST("/hotspot-profiles", CreateHotspotUserProfile)
	control.GET("/hotspot-profiles", ListHotspotUserProfiles)

	// Read-only access
	r.GET("/routers/:id/interfaces", GetInterfaces)
}
