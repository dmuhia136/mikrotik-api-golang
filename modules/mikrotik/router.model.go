package mikrotik

import "gorm.io/gorm"

type Router struct {
	gorm.Model
	Name     string `json:"name"`
	Host     string `json:"host"` // IP or domain
	Port     int    `json:"port"` // 8728 / 8729
	Username string `json:"username"`
	Password string `json:"-"`    // never return
	UseTLS   bool   `json:"use_tls"`
}
