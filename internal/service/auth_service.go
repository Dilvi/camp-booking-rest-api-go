package service

import (
	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
	"github.com/dilvi/camp-booking-rest-api-go/internal/dto"
	"github.com/dilvi/camp-booking-rest-api-go/internal/repository/postgres"
	"github.com/dilvi/camp-booking-rest-api-go/internal/utils"
)

type AuthService struct {
	userRepo *postgres.UserRepository
}

func NewAuthService(userRepo *postgres.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Register(req dto.RegisterRequest) (domain.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "parent",
	}

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return domain.User{}, err
	}

	return createdUser, nil
}