package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/dilvi/camp-booking-rest-api-go/internal/dto"
	"github.com/dilvi/camp-booking-rest-api-go/internal/service"
)

type CampHandler struct {
	campService *service.CampService
}

func NewCampHandler(campService *service.CampService) *CampHandler {
	return &CampHandler{campService: campService}
}

func (h *CampHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	camps, err := h.campService.GetAll()
	if err != nil {
		http.Error(w, "failed to get camps", http.StatusInternalServerError)
		return
	}

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

func (h *CampHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/camps/")
	campID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid camp id", http.StatusBadRequest)
		return
	}

	camp, err := h.campService.GetByID(campID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "camp not found", http.StatusNotFound)
			return
		}
		http.Error(w, "failed to get camp", http.StatusInternalServerError)
		return
	}

	resp := dto.CampResponse{
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
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}