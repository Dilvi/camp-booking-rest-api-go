package domain

import "time"

type Child struct {
	ID        int64
	UserID    int64
	PhotoURL  string
	FirstName string
	LastName  string
	BirthDate time.Time
	Gender    string
	Hobby     string
	Allergy   string
	CreatedAt time.Time
	UpdatedAt time.Time
}