package service

import (
	"errors"
	"week-01-layered-architecture/model"
	"week-01-layered-architecture/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo}
}

func (s *UserService) GetAllUsers() []model.User {
	return s.repo.GetAll()
}

func (s *UserService) CreateUser(user model.User) (model.User, error) {
	if user.Email == "" {
		return model.User{}, errors.New("email is required")
	}
	return s.repo.Create(user), nil
}
