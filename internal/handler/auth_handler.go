package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/dto"
	"github.com/dilvi/camp-booking-rest-api-go/internal/middleware"
	"github.com/dilvi/camp-booking-rest-api-go/internal/respond"
	"github.com/dilvi/camp-booking-rest-api-go/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.authService.Register(req)
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			respond.Error(w, http.StatusConflict, err.Error())
			return
		}
		respond.Error(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	resp := dto.RegisterResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Phone:     user.Phone,
		Email:     user.Email,
		Role:      user.Role,
	}

	respond.JSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.authService.Login(req)
	if err != nil {
		respond.Error(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	resp := dto.LoginResponse{
		Token: token,
	}

	respond.JSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	resp := map[string]interface{}{
		"user_id": claims.UserID,
		"email":   claims.Email,
		"role":    claims.Role,
	}

	respond.JSON(w, http.StatusOK, resp)
}
