package mikrotik

import "mikrotik-api/config"

func LogAction(userID uint, routerID uint, action string, details string) {
	log := AuditLog{
		UserID:   userID,
		RouterID: routerID,
		Action:   action,
		Details:  details,
	}

	config.DB.Create(&log)
}
