package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/dilvi/camp-booking-rest-api-go/internal/dto"
	"github.com/dilvi/camp-booking-rest-api-go/internal/respond"
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
		respond.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	camps, err := h.campService.GetAll()
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get camps")
		return
	}

	resp := make([]dto.CampResponse, 0, len(camps))
	for _, camp := range camps {
		resp = append(resp, dto.CampResponse{
			ID:                camp.ID,
			Title:             camp.Title,
			Location:          camp.Location,
			ImageURL:          camp.ImageURL,
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
	respond.JSON(w, http.StatusOK, resp)
}

func (h *CampHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/camps/")
	campID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid camp id")
		return
	}

	camp, err := h.campService.GetByID(campID)
	if err != nil {
		if err == sql.ErrNoRows {
			respond.Error(w, http.StatusNotFound, "camp not found")
			return
		}
		respond.Error(w, http.StatusInternalServerError, "failed to get camp")
		return
	}

	resp := dto.CampResponse{
		ID:                camp.ID,
		Title:             camp.Title,
		Location:          camp.Location,
		ImageURL:          camp.ImageURL,
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
	respond.JSON(w, http.StatusOK, resp)
}