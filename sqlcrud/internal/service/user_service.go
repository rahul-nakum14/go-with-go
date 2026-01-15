package service

import (
	"github.com/rahulnakum/sqlcrud/internal/model"
	"github.com/rahulnakum/sqlcrud/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(u *model.User) error {
	return s.repo.Create(u)
}

func (s *UserService) Get(id int64) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) Update(u *model.User) error {
	return s.repo.Update(u)
}

func (s *UserService) Delete(id int64) error {
	return s.repo.Delete(id)
}
