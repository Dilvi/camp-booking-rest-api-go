package service

import (
	"github.com/dilvi/camp-booking-rest-api-go/internal/domain"
	"github.com/dilvi/camp-booking-rest-api-go/internal/dto"
	"github.com/dilvi/camp-booking-rest-api-go/internal/repository/postgres"
	"github.com/dilvi/camp-booking-rest-api-go/internal/utils"
)

type ProfileService struct {
	userRepo *postgres.UserRepository
}

func NewProfileService(userRepo *postgres.UserRepository) *ProfileService {
	return &ProfileService{userRepo: userRepo}
}

func (s *ProfileService) GetByUserID(userID int64) (domain.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *ProfileService) Update(userID int64, req dto.UpdateProfileRequest) (domain.User, error) {
	user := domain.User{
		ID:        userID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Email:     req.Email,
		AvatarURL: req.AvatarURL,
	}

	return s.userRepo.UpdateProfile(user)
}

func (s *ProfileService) UpdatePassword(userID int64, req dto.UpdatePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if err := utils.CheckPassword(req.CurrentPassword, user.PasswordHash); err != nil {
		return err
	}

	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, newHash)
}