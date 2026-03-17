package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/dilvi/camp-booking-rest-api-go/internal/dto"
	"github.com/dilvi/camp-booking-rest-api-go/internal/middleware"
	"github.com/dilvi/camp-booking-rest-api-go/internal/service"
)

type FavoriteHandler struct {
	service *service.FavoriteService
}

func NewFavoriteHandler(service *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: service}
}

func (h *FavoriteHandler) Add(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetUserFromContext(r.Context())

	idStr := strings.TrimPrefix(r.URL.Path, "/favorites/")
	campID, _ := strconv.ParseInt(idStr, 10, 64)

	_ = h.service.Add(claims.UserID, campID)

	w.WriteHeader(http.StatusCreated)
}

func (h *FavoriteHandler) Remove(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetUserFromContext(r.Context())

	idStr := strings.TrimPrefix(r.URL.Path, "/favorites/")
	campID, _ := strconv.ParseInt(idStr, 10, 64)

	_ = h.service.Remove(claims.UserID, campID)

	w.WriteHeader(http.StatusNoContent)
}

func (h *FavoriteHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, _ := middleware.GetUserFromContext(r.Context())

	camps, _ := h.service.GetAll(claims.UserID)

	resp := make([]dto.CampResponse, 0, len(camps))
	for _, camp := range camps {
		resp = append(resp, dto.CampResponse{
			ID:                camp.ID,
			Title:             camp.Title,
			Location:          camp.Location,
			PricePerDay:       camp.PricePerDay,
			BookedCount:       camp.BookedCount,
			Description:       camp.Description,
			ShiftDurationDays: camp.ShiftDurationDays,
			AgeMin:            camp.AgeMin,
			AgeMax:            camp.AgeMax,
			CampType:          camp.CampType,
			FoodType:          camp.FoodType,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}