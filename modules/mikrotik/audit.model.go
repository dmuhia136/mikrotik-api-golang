package mikrotik

import "gorm.io/gorm"

type AuditLog struct {
	gorm.Model
	UserID   uint
	RouterID uint
	Action   string
	Details  string
}
