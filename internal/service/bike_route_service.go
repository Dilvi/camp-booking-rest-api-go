package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
	"github.com/dilvi/camp-booking-rest-api-go/internal/gigachat"
	"github.com/dilvi/camp-booking-rest-api-go/internal/repository/postgres"
)

type BikeRouteService struct {
	repo           *postgres.BikeRouteRepository
	bookingRepo    *postgres.BookingRepository
	gigachatClient *gigachat.Client
}

var (
	ErrBikeRouteBookingRequired  = errors.New("bike route analysis requires booking")
	ErrBikeRouteNotInBookingCamp = errors.New("bike route does not belong to booking camp")
	ErrBikeRouteInactiveBooking  = errors.New("booking is not active")
)

func NewBikeRouteService(repo *postgres.BikeRouteRepository, bookingRepo *postgres.BookingRepository, gigachatClient *gigachat.Client) *BikeRouteService {
	return &BikeRouteService{
		repo:           repo,
		bookingRepo:    bookingRepo,
		gigachatClient: gigachatClient,
	}
}

func (s *BikeRouteService) GetAll(campID *int64) ([]domain.BikeRoute, error) {
	return s.repo.GetAll(campID)
}

func (s *BikeRouteService) GetByID(id int64) (domain.BikeRoute, error) {
	return s.repo.GetByID(id)
}

func (s *BikeRouteService) Analyze(ctx context.Context, userID, id int64) (domain.BikeRoute, string, error) {
	route, err := s.repo.GetByID(id)
	if err != nil {
		return domain.BikeRoute{}, "", err
	}

	hasBooking, err := s.bookingRepo.ExistsActiveByUserAndCamp(userID, route.CampID)
	if err != nil {
		return domain.BikeRoute{}, "", err
	}
	if !hasBooking {
		return domain.BikeRoute{}, "", ErrBikeRouteBookingRequired
	}

	prompt := buildBikeRouteAnalysisPrompt(route)
	recommendations, err := s.gigachatClient.Complete(ctx, prompt)
	if err != nil {
		return domain.BikeRoute{}, "", err
	}

	return route, recommendations, nil
}

func (s *BikeRouteService) GetAllForBooking(userID, bookingID int64) ([]domain.BikeRoute, error) {
	booking, err := s.bookingRepo.GetByIDForUser(bookingID, userID)
	if err != nil {
		return nil, err
	}
	if !isActiveBookingStatus(booking.Status) {
		return nil, ErrBikeRouteInactiveBooking
	}

	return s.repo.GetAll(&booking.CampID)
}

func (s *BikeRouteService) AnalyzeForBooking(ctx context.Context, userID, bookingID, routeID int64) (domain.BikeRoute, string, error) {
	booking, err := s.bookingRepo.GetByIDForUser(bookingID, userID)
	if err != nil {
		return domain.BikeRoute{}, "", err
	}
	if !isActiveBookingStatus(booking.Status) {
		return domain.BikeRoute{}, "", ErrBikeRouteInactiveBooking
	}

	route, err := s.repo.GetByID(routeID)
	if err != nil {
		return domain.BikeRoute{}, "", err
	}
	if route.CampID != booking.CampID {
		return domain.BikeRoute{}, "", ErrBikeRouteNotInBookingCamp
	}

	prompt := buildBikeRouteAnalysisPrompt(route)
	recommendations, err := s.gigachatClient.Complete(ctx, prompt)
	if err != nil {
		return domain.BikeRoute{}, "", err
	}

	return route, recommendations, nil
}

func isActiveBookingStatus(status string) bool {
	return status == "pending" || status == "confirmed"
}

func buildBikeRouteAnalysisPrompt(route domain.BikeRoute) string {
	return fmt.Sprintf(`Проанализируй веломаршрут для семьи с ребенком, которая отдыхает в детском лагере.

Данные маршрута:
- Название: %s
- Локация: %s
- Дистанция: %.1f км
- Примерная длительность: %d минут
- Сложность: %s
- Тип маршрута: %s
- Набор высоты: %d м
- Описание: %s
- Точки маршрута или геоданные: %s
- Дополнительный контекст: %s

Верни практичные рекомендации на русском языке:
1. Кому подойдет маршрут.
2. Что взять с собой.
3. На что обратить внимание по безопасности.
4. Лучшее время для поездки.
5. Краткий итог в 1-2 предложениях.`,
		route.Title,
		route.Location,
		route.DistanceKM,
		route.DurationMinutes,
		route.Difficulty,
		route.RouteType,
		route.ElevationGainM,
		route.Description,
		route.RoutePoints,
		route.RecommendationsContext,
	)
}
