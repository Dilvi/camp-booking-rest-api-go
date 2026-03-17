package domain

import "time"

type User struct {
	ID           int64
	FirstName    string
	LastName     string
	Phone        string
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}