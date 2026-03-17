package handler

import (
	"net/http"

	"github.com/dilvi/camp-booking-rest-api-go/internal/respond"
)

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"status": "ok",
	}
	respond.JSON(w, http.StatusOK, resp)
}