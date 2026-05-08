package service

import (
	"database/sql"
	"errors"

	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
	"github.com/dilvi/camp-booking-rest-api-go/internal/repository/postgres"
)

type BookingService struct {
	repo      *postgres.BookingRepository
	childRepo *postgres.ChildRepository
	campRepo  *postgres.CampRepository
}

func NewBookingService(
	repo *postgres.BookingRepository,
	childRepo *postgres.ChildRepository,
	campRepo *postgres.CampRepository,
) *BookingService {
	return &BookingService{
		repo:      repo,
		childRepo: childRepo,
		campRepo:  campRepo,
	}
}

var (
	ErrBookingChildDoesNotBelongToUser = errors.New("child does not belong to user")
	ErrBookingCampNotFound             = errors.New("camp not found")
)

func (s *BookingService) Create(userID, childID, campID int64) (domain.Booking, error) {
	if childID <= 0 || campID <= 0 {
		return domain.Booking{}, errors.New("child_id and camp_id are required")
	}

	children, err := s.childRepo.GetAllByUserID(userID)
	if err != nil {
		return domain.Booking{}, err
	}

	valid := false
	for _, c := range children {
		if c.ID == childID {
			valid = true
			break
		}
	}

	if !valid {
		return domain.Booking{}, ErrBookingChildDoesNotBelongToUser
	}

	_, err = s.campRepo.GetByID(campID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Booking{}, ErrBookingCampNotFound
		}
		return domain.Booking{}, err
	}

	return s.repo.Create(domain.Booking{
		UserID:  userID,
		ChildID: childID,
		CampID:  campID,
	})
}

func (s *BookingService) GetAll(userID int64) ([]domain.Booking, error) {
	return s.repo.GetAllByUserID(userID)
}

func (s *BookingService) GetByID(userID, bookingID int64) (domain.Booking, error) {
	return s.repo.GetByIDForUser(bookingID, userID)
}
