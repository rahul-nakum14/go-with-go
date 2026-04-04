package service

import (
	"log"

	"github.com/rahul-nakum14/go-user-crud/internal/models"
	"github.com/rahul-nakum14/go-user-crud/internal/repository"
)

type UserService struct {
	UserRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
	}
}

func (s *UserService) GetAllUsers()([]models.User, error) {
	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return nil, err
	}
	log.Printf("Fetched users: %v", users)
	return users, nil
}

func (s *UserService) CreateUser(username, email, password string) (int64, error) {
	user := models.User{
		Username: username,
		Email:    email,
		Password: password,
	}
	id, err := s.UserRepo.CreateUser(user)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return 0, err
	}
	log.Printf("Created user with ID: %d", id)
	return id, nil
}