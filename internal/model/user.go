package model

import "time"

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:100"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"default:user"`
	CreatedAt    time.Time
}
