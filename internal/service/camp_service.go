package service

import (
	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
	"github.com/dilvi/camp-booking-rest-api-go/internal/repository/postgres"
)

type CampService struct {
	campRepo *postgres.CampRepository
}

func NewCampService(campRepo *postgres.CampRepository) *CampService {
	return &CampService{campRepo: campRepo}
}

func (s *CampService) GetAll() ([]domain.Camp, error) {
	return s.campRepo.GetAll()
}

func (s *CampService) GetByID(id int64) (domain.Camp, error) {
	return s.campRepo.GetByID(id)
}