package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/dilvi/camp-booking-rest-api-go/internal/dto"
	"github.com/dilvi/camp-booking-rest-api-go/internal/middleware"
	"github.com/dilvi/camp-booking-rest-api-go/internal/respond"
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
		respond.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreateChildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	child, err := h.childService.Create(claims.UserID, req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "failed to create child")
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
	respond.JSON(w, http.StatusOK, resp)
}

func (h *ChildHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	children, err := h.childService.GetAllByUserID(claims.UserID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get children")
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
	respond.JSON(w, http.StatusOK, resp)
}

func (h *ChildHandler) Update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		respond.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/children/")
	childID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid child id")
		return
	}

	var req dto.UpdateChildRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	child, err := h.childService.Update(claims.UserID, childID, req)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "failed to update child")
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
	respond.JSON(w, http.StatusOK, resp)
}

func (h *ChildHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		respond.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/children/")
	childID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid child id")
		return
	}

	err = h.childService.Delete(claims.UserID, childID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to delete child")
		return
	}

	respond.JSON(w, http.StatusOK, map[string]string{
		"message": "child deleted successfully",
	})
}