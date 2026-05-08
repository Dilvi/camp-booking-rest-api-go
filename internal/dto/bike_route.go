package dto

type BikeRouteResponse struct {
	ID              int64   `json:"id"`
	CampID          int64   `json:"camp_id"`
	Title           string  `json:"title"`
	Location        string  `json:"location"`
	MapURL          string  `json:"map_url"`
	RoutePoints     string  `json:"route_points"`
	DistanceKM      float64 `json:"distance_km"`
	DurationMinutes int     `json:"duration_minutes"`
	Difficulty      string  `json:"difficulty"`
	RouteType       string  `json:"route_type"`
	ElevationGainM  int     `json:"elevation_gain_m"`
	Description     string  `json:"description"`
}

type BikeRouteAIAnalysisResponse struct {
	Route           BikeRouteResponse `json:"route"`
	Recommendations string            `json:"recommendations"`
}
