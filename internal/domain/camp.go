package domain

import "time"

type Camp struct {
	ID                int64
	Title             string
	Location          string
	PricePerDay       int
	BookedCount       int
	Description       string
	ShiftDurationDays int
	AgeMin            int
	AgeMax            int
	CampType          string
	FoodType          string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}