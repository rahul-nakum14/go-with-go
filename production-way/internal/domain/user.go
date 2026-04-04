// Package domain defines the core business entities and repository contracts.
// No framework imports — pure Go types only.
package domain

import (
	"context"
	"errors"
	"time"
)

// ── Sentinel errors ────────────────────────────────────────────────────────────

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrEmailAlreadyTaken = errors.New("email already taken")
	ErrInvalidID         = errors.New("invalid user id")
)

// ── Entity ─────────────────────────────────────────────────────────────────────

// User is the core domain entity. Keep it framework-agnostic.
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Role is a typed string to prevent stringly-typed bugs.
type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

func (r Role) Valid() bool {
	return r == RoleAdmin || r == RoleUser
}

// ── Value objects / DTOs ───────────────────────────────────────────────────────

// CreateUserInput holds validated data for creating a user.
type CreateUserInput struct {
	Name  string `json:"name"  validate:"required,min=2,max=100"`
	Email string `json:"email" validate:"required,email"`
	Role  Role   `json:"role"  validate:"required,oneof=admin user"`
}

// UpdateUserInput holds validated data for updating a user.
// Pointer fields = optional / partial update semantics.
type UpdateUserInput struct {
	Name  *string `json:"name"  validate:"omitempty,min=2,max=100"`
	Email *string `json:"email" validate:"omitempty,email"`
	Role  *Role   `json:"role"  validate:"omitempty,oneof=admin user"`
}

// ListFilter controls pagination and filtering for user lists.
type ListFilter struct {
	Page    int    `json:"page"    validate:"min=1"`
	PerPage int    `json:"per_page" validate:"min=1,max=100"`
	Search  string `json:"search"`
}

func (f *ListFilter) Offset() int {
	return (f.Page - 1) * f.PerPage
}

// ── Repository contract ────────────────────────────────────────────────────────

// UserRepository is the persistence abstraction.
// The service layer depends on this interface — NOT on a concrete DB type.
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	List(ctx context.Context, filter ListFilter) ([]*User, int64, error)
	Create(ctx context.Context, input CreateUserInput) (*User, error)
	Update(ctx context.Context, id int64, input UpdateUserInput) (*User, error)
	Delete(ctx context.Context, id int64) error
}
