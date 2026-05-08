package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
	"github.com/dilvi/camp-booking-rest-api-go/internal/dto"
	"github.com/dilvi/camp-booking-rest-api-go/internal/middleware"
	"github.com/dilvi/camp-booking-rest-api-go/internal/respond"
	"github.com/dilvi/camp-booking-rest-api-go/internal/service"
)

type BookingHandler struct {
	service          *service.BookingService
	bikeRouteService *service.BikeRouteService
}

func NewBookingHandler(service *service.BookingService, bikeRouteService *service.BikeRouteService) *BookingHandler {
	return &BookingHandler{
		service:          service,
		bikeRouteService: bikeRouteService,
	}
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req dto.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid body")
		return
	}

	booking, err := h.service.Create(claims.UserID, req.ChildID, req.CampID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrBookingChildDoesNotBelongToUser):
			respond.Error(w, http.StatusForbidden, err.Error())
		case errors.Is(err, service.ErrBookingCampNotFound):
			respond.Error(w, http.StatusNotFound, err.Error())
		default:
			respond.Error(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	respond.JSON(w, http.StatusCreated, bookingToResponse(booking))
}

func (h *BookingHandler) List(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bookings, err := h.service.GetAll(claims.UserID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get bookings")
		return
	}

	resp := make([]dto.BookingResponse, 0, len(bookings))
	for _, booking := range bookings {
		resp = append(resp, bookingToResponse(booking))
	}

	respond.JSON(w, http.StatusOK, resp)
}

func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bookingID, err := bookingIDFromPath(r.URL.Path)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid booking id")
		return
	}

	booking, err := h.service.GetByID(claims.UserID, bookingID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respond.Error(w, http.StatusNotFound, "booking not found")
			return
		}
		respond.Error(w, http.StatusInternalServerError, "failed to get booking")
		return
	}

	respond.JSON(w, http.StatusOK, bookingToResponse(booking))
}

func (h *BookingHandler) ListBikeRoutes(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bookingID, err := bookingIDFromBikeRoutesPath(r.URL.Path)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid booking id")
		return
	}

	routes, err := h.bikeRouteService.GetAllForBooking(claims.UserID, bookingID)
	if err != nil {
		writeBikeRouteBookingError(w, err)
		return
	}

	resp := make([]dto.BikeRouteResponse, 0, len(routes))
	for _, route := range routes {
		resp = append(resp, bikeRouteToResponse(route))
	}

	respond.JSON(w, http.StatusOK, resp)
}

func (h *BookingHandler) BikeRouteAIAnalysis(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bookingID, routeID, err := bookingAndRouteIDsFromAnalysisPath(r.URL.Path)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid booking or bike route id")
		return
	}

	route, recommendations, err := h.bikeRouteService.AnalyzeForBooking(r.Context(), claims.UserID, bookingID, routeID)
	if err != nil {
		writeBikeRouteBookingError(w, err)
		return
	}

	respond.JSON(w, http.StatusOK, dto.BikeRouteAIAnalysisResponse{
		Route:           bikeRouteToResponse(route),
		Recommendations: recommendations,
	})
}

func bookingToResponse(booking domain.Booking) dto.BookingResponse {
	return dto.BookingResponse{
		ID:        booking.ID,
		ChildID:   booking.ChildID,
		CampID:    booking.CampID,
		Status:    booking.Status,
		CreatedAt: formatTime(booking.CreatedAt),
		UpdatedAt: formatTime(booking.UpdatedAt),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(time.RFC3339)
}

func bookingIDFromPath(path string) (int64, error) {
	idStr := strings.TrimPrefix(path, "/bookings/")
	if idStr == "" || strings.Contains(idStr, "/") {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseInt(idStr, 10, 64)
}

func bookingIDFromBikeRoutesPath(path string) (int64, error) {
	parts := strings.Split(strings.TrimPrefix(path, "/bookings/"), "/")
	if len(parts) != 2 || parts[1] != "bike-routes" {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseInt(parts[0], 10, 64)
}

func bookingAndRouteIDsFromAnalysisPath(path string) (int64, int64, error) {
	parts := strings.Split(strings.TrimPrefix(path, "/bookings/"), "/")
	if len(parts) != 4 || parts[1] != "bike-routes" || parts[3] != "ai-analysis" {
		return 0, 0, strconv.ErrSyntax
	}

	bookingID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	routeID, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return bookingID, routeID, nil
}

func writeBikeRouteBookingError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		respond.Error(w, http.StatusNotFound, "booking or bike route not found")
	case errors.Is(err, service.ErrBikeRouteInactiveBooking):
		respond.Error(w, http.StatusForbidden, "booking is not active")
	case errors.Is(err, service.ErrBikeRouteNotInBookingCamp):
		respond.Error(w, http.StatusForbidden, "bike route does not belong to booking camp")
	case errors.Is(err, service.ErrBikeRouteBookingRequired):
		respond.Error(w, http.StatusForbidden, "booking required for bike route analysis")
	default:
		log.Printf("bike route booking request failed: %v", err)
		respond.Error(w, http.StatusBadGateway, "failed to process bike route request")
	}
}
