package dto

type CreateBookingRequest struct {
	ChildID int64 `json:"child_id"`
	CampID  int64 `json:"camp_id"`
}

type BookingResponse struct {
	ID      int64  `json:"id"`
	ChildID int64  `json:"child_id"`
	CampID  int64  `json:"camp_id"`
	Status  string `json:"status"`
}