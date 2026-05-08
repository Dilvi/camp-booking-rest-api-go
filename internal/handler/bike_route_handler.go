package handler

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
	"github.com/dilvi/camp-booking-rest-api-go/internal/dto"
	"github.com/dilvi/camp-booking-rest-api-go/internal/middleware"
	"github.com/dilvi/camp-booking-rest-api-go/internal/respond"
	"github.com/dilvi/camp-booking-rest-api-go/internal/service"
)

type BikeRouteHandler struct {
	service *service.BikeRouteService
}

func NewBikeRouteHandler(service *service.BikeRouteService) *BikeRouteHandler {
	return &BikeRouteHandler{service: service}
}

func (h *BikeRouteHandler) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var campID *int64
	if campIDParam := r.URL.Query().Get("camp_id"); campIDParam != "" {
		parsedCampID, err := strconv.ParseInt(campIDParam, 10, 64)
		if err != nil {
			respond.Error(w, http.StatusBadRequest, "invalid camp id")
			return
		}
		campID = &parsedCampID
	}

	routes, err := h.service.GetAll(campID)
	if err != nil {
		respond.Error(w, http.StatusInternalServerError, "failed to get bike routes")
		return
	}

	resp := make([]dto.BikeRouteResponse, 0, len(routes))
	for _, route := range routes {
		resp = append(resp, bikeRouteToResponse(route))
	}

	respond.JSON(w, http.StatusOK, resp)
}

func (h *BikeRouteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respond.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	routeID, err := bikeRouteIDFromPath(r.URL.Path)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid bike route id")
		return
	}

	route, err := h.service.GetByID(routeID)
	if err != nil {
		if err == sql.ErrNoRows {
			respond.Error(w, http.StatusNotFound, "bike route not found")
			return
		}
		respond.Error(w, http.StatusInternalServerError, "failed to get bike route")
		return
	}

	respond.JSON(w, http.StatusOK, bikeRouteToResponse(route))
}

func (h *BikeRouteHandler) AIAnalysis(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respond.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		respond.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	routeID, err := bikeRouteIDFromAnalysisPath(r.URL.Path)
	if err != nil {
		respond.Error(w, http.StatusBadRequest, "invalid bike route id")
		return
	}

	route, recommendations, err := h.service.Analyze(r.Context(), claims.UserID, routeID)
	if err != nil {
		if err == sql.ErrNoRows {
			respond.Error(w, http.StatusNotFound, "bike route not found")
			return
		}
		if errors.Is(err, service.ErrBikeRouteBookingRequired) {
			respond.Error(w, http.StatusForbidden, "booking required for bike route analysis")
			return
		}
		log.Printf("bike route analysis failed: %v", err)
		respond.Error(w, http.StatusBadGateway, "failed to analyze bike route")
		return
	}

	respond.JSON(w, http.StatusOK, dto.BikeRouteAIAnalysisResponse{
		Route:           bikeRouteToResponse(route),
		Recommendations: recommendations,
	})
}

func bikeRouteIDFromPath(path string) (int64, error) {
	idStr := strings.TrimPrefix(path, "/bike-routes/")
	if idStr == "" || strings.Contains(idStr, "/") {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseInt(idStr, 10, 64)
}

func bikeRouteIDFromAnalysisPath(path string) (int64, error) {
	idStr := strings.TrimPrefix(path, "/bike-routes/")
	idStr = strings.TrimSuffix(idStr, "/ai-analysis")
	if idStr == "" || strings.Contains(idStr, "/") {
		return 0, strconv.ErrSyntax
	}
	return strconv.ParseInt(idStr, 10, 64)
}

func bikeRouteToResponse(route domain.BikeRoute) dto.BikeRouteResponse {
	return dto.BikeRouteResponse{
		ID:              route.ID,
		CampID:          route.CampID,
		Title:           route.Title,
		Location:        route.Location,
		MapURL:          route.MapURL,
		RoutePoints:     route.RoutePoints,
		DistanceKM:      route.DistanceKM,
		DurationMinutes: route.DurationMinutes,
		Difficulty:      route.Difficulty,
		RouteType:       route.RouteType,
		ElevationGainM:  route.ElevationGainM,
		Description:     route.Description,
	}
}
