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

type ChildHandler struct {
	childService *service.ChildService
}

func NewChildHandler(childService *service.ChildService) *ChildHandler {
	return &ChildHandler{childService: childService}
}

func (h *ChildHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req dto.CreateChildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	child, err := h.childService.Create(claims.UserID, req)
	if err != nil {
		http.Error(w, "failed to create child", http.StatusBadRequest)
		return
	}

	resp := dto.ChildResponse{
		ID:        child.ID,
		PhotoURL:  child.PhotoURL,
		FirstName: child.FirstName,
		LastName:  child.LastName,
		BirthDate: child.BirthDate.Format("2006-01-02"),
		Gender:    child.Gender,
		Hobby:     child.Hobby,
		Allergy:   child.Allergy,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *ChildHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	children, err := h.childService.GetAllByUserID(claims.UserID)
	if err != nil {
		http.Error(w, "failed to get children", http.StatusInternalServerError)
		return
	}

	resp := make([]dto.ChildResponse, 0, len(children))
	for _, child := range children {
		resp = append(resp, dto.ChildResponse{
			ID:        child.ID,
			PhotoURL:  child.PhotoURL,
			FirstName: child.FirstName,
			LastName:  child.LastName,
			BirthDate: child.BirthDate.Format("2006-01-02"),
			Gender:    child.Gender,
			Hobby:     child.Hobby,
			Allergy:   child.Allergy,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (h *ChildHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/children/")
	childID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid child id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateChildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	child, err := h.childService.Update(claims.UserID, childID, req)
	if err != nil {
		http.Error(w, "failed to update child", http.StatusBadRequest)
		return
	}

	resp := dto.ChildResponse{
		ID:        child.ID,
		PhotoURL:  child.PhotoURL,
		FirstName: child.FirstName,
		LastName:  child.LastName,
		BirthDate: child.BirthDate.Format("2006-01-02"),
		Gender:    child.Gender,
		Hobby:     child.Hobby,
		Allergy:   child.Allergy,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}