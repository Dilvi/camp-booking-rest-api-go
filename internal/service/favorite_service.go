package service

import (
	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
	"github.com/dilvi/camp-booking-rest-api-go/internal/repository/postgres"
)

type FavoriteService struct {
	repo *postgres.FavoriteRepository
}

func NewFavoriteService(repo *postgres.FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo}
}

func (s *FavoriteService) Add(userID, campID int64) error {
	return s.repo.Add(userID, campID)
}

func (s *FavoriteService) Remove(userID, campID int64) error {
	return s.repo.Remove(userID, campID)
}

func (s *FavoriteService) GetAll(userID int64) ([]domain.Camp, error) {
	return s.repo.GetAll(userID)
}