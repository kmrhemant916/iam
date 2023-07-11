package service

import (
	"github.com/kmrhemant916/iam/entities"
	"github.com/kmrhemant916/iam/repositories"
)

type SignupService interface {
	CreateRootAccount(user *entities.User, organization *entities.Organization) error
}

type signupService struct {
	signupRepository repositories.SignupRepository
}

func NewSignupService(signupRepository repositories.SignupRepository) *signupService {
	return &signupService{
		signupRepository,
	}
}

func (s *signupService) CreateRootAccount(user *entities.User, organization *entities.Organization) error {
	err := s.signupRepository.CreateRootAccount(user, organization)
	return err
}