package dto

type CampResponse struct {
	ID                int64  `json:"id"`
	Title             string `json:"title"`
	Location          string `json:"location"`
	PricePerDay       int    `json:"price_per_day"`
	BookedCount       int    `json:"booked_count"`
	Description       string `json:"description"`
	ShiftDurationDays int    `json:"shift_duration_days"`
	AgeMin            int    `json:"age_min"`
	AgeMax            int    `json:"age_max"`
	CampType          string `json:"camp_type"`
	FoodType          string `json:"food_type"`
}