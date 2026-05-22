package entity

import "time"

type User struct {
	ID           string
	Email        string
	Name         string
	AvatarColor  string
	Role         string
	Disabled     bool
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

