package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/dto"
	"github.com/dilvi/camp-booking-rest-api-go/internal/middleware"
	"github.com/dilvi/camp-booking-rest-api-go/internal/respond"
	"github.com/dilvi/camp-booking-rest-api-go/internal/service"
)

type BookingHandler struct {
	service *service.BookingService
}

func NewBookingHandler(service *service.BookingService) *BookingHandler {
	return &BookingHandler{service: service}
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetUserFromContext(r.Context())

	var req dto.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	booking, err := h.service.Create(claims.UserID, req.ChildID, req.CampID)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := dto.BookingResponse{
		ID:      booking.ID,
		ChildID: booking.ChildID,
		CampID:  booking.CampID,
		Status:  booking.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	respond.JSON(w, http.StatusOK, resp)
}

func (h *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetUserFromContext(r.Context())

	bookings, _ := h.service.GetAll(claims.UserID)

	resp := make([]dto.BookingResponse, 0, len(bookings))
	for _, b := range bookings {
		resp = append(resp, dto.BookingResponse{
			ID:      b.ID,
			ChildID: b.ChildID,
			CampID:  b.CampID,
			Status:  b.Status,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	respond.JSON(w, http.StatusOK, resp)
}