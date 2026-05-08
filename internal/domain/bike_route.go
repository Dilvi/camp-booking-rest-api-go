package domain

import "time"

type BikeRoute struct {
	ID                     int64
	CampID                 int64
	Title                  string
	Location               string
	MapURL                 string
	RoutePoints            string
	DistanceKM             float64
	DurationMinutes        int
	Difficulty             string
	RouteType              string
	ElevationGainM         int
	Description            string
	RecommendationsContext string
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
