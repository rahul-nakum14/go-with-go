// Package service implements business logic. It depends only on the domain
// interfaces — not on HTTP, SQL, or any infrastructure concern.
package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/yourorg/go-user-crud/internal/domain"
)

// UserService encapsulates all user-related business logic.
type UserService struct {
	repo domain.UserRepository
	log  *slog.Logger
}

// NewUserService constructs a UserService with its required dependencies.
func NewUserService(repo domain.UserRepository, log *slog.Logger) *UserService {
	return &UserService{repo: repo, log: log}
}

// GetByID retrieves a single user by ID.
func (s *UserService) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("GetByID: %w", err)
	}
	return user, nil
}

// List returns a paginated list of users with optional search.
func (s *UserService) List(ctx context.Context, filter domain.ListFilter) ([]*domain.User, int64, error) {
	users, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("List: %w", err)
	}
	return users, total, nil
}

// Create validates business rules then persists a new user.
func (s *UserService) Create(ctx context.Context, input domain.CreateUserInput) (*domain.User, error) {
	// Business rule: email must be unique.
	existing, err := s.repo.GetByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
		return nil, fmt.Errorf("Create: checking email: %w", err)
	}
	if existing != nil {
		return nil, domain.ErrEmailAlreadyTaken
	}

	user, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("Create: %w", err)
	}

	s.log.InfoContext(ctx, "user created", "user_id", user.ID, "email", user.Email)
	return user, nil
}

// Update applies partial updates to an existing user.
func (s *UserService) Update(ctx context.Context, id int64, input domain.UpdateUserInput) (*domain.User, error) {
	// Ensure the user exists first.
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("Update: fetch check: %w", err)
	}

	// If email is changing, enforce uniqueness.
	if input.Email != nil {
		existing, err := s.repo.GetByEmail(ctx, *input.Email)
		if err != nil && !errors.Is(err, domain.ErrUserNotFound) {
			return nil, fmt.Errorf("Update: checking email: %w", err)
		}
		if existing != nil && existing.ID != id {
			return nil, domain.ErrEmailAlreadyTaken
		}
	}

	user, err := s.repo.Update(ctx, id, input)
	if err != nil {
		return nil, fmt.Errorf("Update: %w", err)
	}

	s.log.InfoContext(ctx, "user updated", "user_id", id)
	return user, nil
}

// Delete removes a user permanently.
func (s *UserService) Delete(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("Delete: fetch check: %w", err)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("Delete: %w", err)
	}

	s.log.InfoContext(ctx, "user deleted", "user_id", id)
	return nil
}
