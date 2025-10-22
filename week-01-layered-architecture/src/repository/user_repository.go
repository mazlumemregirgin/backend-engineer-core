package repository

import "week-01-layered-architecture/model"

type UserRepository struct {
	users []model.User
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		users: []model.User{
			{ID: 1, Name: "Mazlum", Email: "mazlum@example.com"},
			{ID: 2, Name: "Emre", Email: "emre@example.com"},
		},
	}
}

func (r *UserRepository) GetAll() []model.User {
	return r.users
}

func (r *UserRepository) Create(user model.User) model.User {
	user.ID = len(r.users) + 1
	r.users = append(r.users, user)
	return user
}
