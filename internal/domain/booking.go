package domain

import "time"

type Booking struct {
	ID        int64
	UserID    int64
	ChildID   int64
	CampID    int64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}