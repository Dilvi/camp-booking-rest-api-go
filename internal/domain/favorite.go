package domain

import "time"

type Favorite struct {
	ID        int64
	UserID    int64
	CampID    int64
	CreatedAt time.Time
}