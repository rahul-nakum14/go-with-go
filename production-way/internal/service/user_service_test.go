package service_test

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/yourorg/go-user-crud/internal/domain"
	"github.com/yourorg/go-user-crud/internal/service"
)

// ── Mock repository ───────────────────────────────────────────────────────────

type mockRepo struct {
	getByIDFn    func(ctx context.Context, id int64) (*domain.User, error)
	getByEmailFn func(ctx context.Context, email string) (*domain.User, error)
	listFn       func(ctx context.Context, f domain.ListFilter) ([]*domain.User, int64, error)
	createFn     func(ctx context.Context, in domain.CreateUserInput) (*domain.User, error)
	updateFn     func(ctx context.Context, id int64, in domain.UpdateUserInput) (*domain.User, error)
	deleteFn     func(ctx context.Context, id int64) error
}

func (m *mockRepo) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return m.getByEmailFn(ctx, email)
}
func (m *mockRepo) List(ctx context.Context, f domain.ListFilter) ([]*domain.User, int64, error) {
	return m.listFn(ctx, f)
}
func (m *mockRepo) Create(ctx context.Context, in domain.CreateUserInput) (*domain.User, error) {
	return m.createFn(ctx, in)
}
func (m *mockRepo) Update(ctx context.Context, id int64, in domain.UpdateUserInput) (*domain.User, error) {
	return m.updateFn(ctx, id, in)
}
func (m *mockRepo) Delete(ctx context.Context, id int64) error { return m.deleteFn(ctx, id) }

// ── Helpers ───────────────────────────────────────────────────────────────────

func newSvc(repo domain.UserRepository) *service.UserService {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return service.NewUserService(repo, log)
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestCreate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     domain.CreateUserInput
		repoSetup func(*mockRepo)
		wantErr   error
	}{
		{
			name:  "success",
			input: domain.CreateUserInput{Name: "Alice", Email: "alice@example.com", Role: domain.RoleUser},
			repoSetup: func(r *mockRepo) {
				r.getByEmailFn = func(_ context.Context, _ string) (*domain.User, error) {
					return nil, domain.ErrUserNotFound
				}
				r.createFn = func(_ context.Context, in domain.CreateUserInput) (*domain.User, error) {
					return &domain.User{ID: 1, Name: in.Name, Email: in.Email, Role: in.Role}, nil
				}
			},
		},
		{
			name:  "duplicate email",
			input: domain.CreateUserInput{Name: "Bob", Email: "taken@example.com", Role: domain.RoleUser},
			repoSetup: func(r *mockRepo) {
				r.getByEmailFn = func(_ context.Context, _ string) (*domain.User, error) {
					return &domain.User{ID: 99}, nil // already exists
				}
			},
			wantErr: domain.ErrEmailAlreadyTaken,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			repo := &mockRepo{}
			tt.repoSetup(repo)
			svc := newSvc(repo)

			_, err := svc.Create(context.Background(), tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestDelete_NotFound(t *testing.T) {
	t.Parallel()

	repo := &mockRepo{
		getByIDFn: func(_ context.Context, _ int64) (*domain.User, error) {
			return nil, domain.ErrUserNotFound
		},
	}

	svc := newSvc(repo)
	err := svc.Delete(context.Background(), 999)
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
