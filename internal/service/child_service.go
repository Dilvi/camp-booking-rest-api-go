package service

import (
	"time"

	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
	"github.com/dilvi/camp-booking-rest-api-go/internal/dto"
	"github.com/dilvi/camp-booking-rest-api-go/internal/repository/postgres"
)

type ChildService struct {
	childRepo *postgres.ChildRepository
}

func NewChildService(childRepo *postgres.ChildRepository) *ChildService {
	return &ChildService{childRepo: childRepo}
}

func (s *ChildService) Create(userID int64, req dto.CreateChildRequest) (domain.Child, error) {
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return domain.Child{}, err
	}

	child := domain.Child{
		UserID:    userID,
		PhotoURL:  req.PhotoURL,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: birthDate,
		Gender:    req.Gender,
		Hobby:     req.Hobby,
		Allergy:   req.Allergy,
	}

	return s.childRepo.Create(child)
}

func (s *ChildService) GetAllByUserID(userID int64) ([]domain.Child, error) {
	return s.childRepo.GetAllByUserID(userID)
}

func (s *ChildService) Update(userID, childID int64, req dto.UpdateChildRequest) (domain.Child, error) {
	birthDate, err := time.Parse("2006-01-02", req.BirthDate)
	if err != nil {
		return domain.Child{}, err
	}

	child := domain.Child{
		ID:        childID,
		UserID:    userID,
		PhotoURL:  req.PhotoURL,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		BirthDate: birthDate,
		Gender:    req.Gender,
		Hobby:     req.Hobby,
		Allergy:   req.Allergy,
	}

	return s.childRepo.Update(child)
}