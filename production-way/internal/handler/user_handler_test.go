package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/yourorg/go-user-crud/internal/domain"
	"github.com/yourorg/go-user-crud/internal/handler"
	"github.com/yourorg/go-user-crud/pkg/validator"
)

// ── Mock service ──────────────────────────────────────────────────────────────

type mockSvc struct {
	getByIDFn func(ctx context.Context, id int64) (*domain.User, error)
	listFn    func(ctx context.Context, f domain.ListFilter) ([]*domain.User, int64, error)
	createFn  func(ctx context.Context, in domain.CreateUserInput) (*domain.User, error)
	updateFn  func(ctx context.Context, id int64, in domain.UpdateUserInput) (*domain.User, error)
	deleteFn  func(ctx context.Context, id int64) error
}

func (m *mockSvc) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockSvc) List(ctx context.Context, f domain.ListFilter) ([]*domain.User, int64, error) {
	return m.listFn(ctx, f)
}
func (m *mockSvc) Create(ctx context.Context, in domain.CreateUserInput) (*domain.User, error) {
	return m.createFn(ctx, in)
}
func (m *mockSvc) Update(ctx context.Context, id int64, in domain.UpdateUserInput) (*domain.User, error) {
	return m.updateFn(ctx, id, in)
}
func (m *mockSvc) Delete(ctx context.Context, id int64) error { return m.deleteFn(ctx, id) }

// ── Helpers ───────────────────────────────────────────────────────────────────

func newHandler(svc *mockSvc) *handler.UserHandler {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	return handler.NewUserHandler(svc, validator.New(), log)
}

func toJSON(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return bytes.NewBuffer(b)
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func TestCreate_Handler(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       any
		svcSetup   func(*mockSvc)
		wantStatus int
	}{
		{
			name: "created",
			body: domain.CreateUserInput{Name: "Alice", Email: "alice@example.com", Role: domain.RoleUser},
			svcSetup: func(s *mockSvc) {
				s.createFn = func(_ context.Context, in domain.CreateUserInput) (*domain.User, error) {
					return &domain.User{ID: 1, Name: in.Name, Email: in.Email, Role: in.Role}, nil
				}
			},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "bad json",
			body:       "not-json",
			svcSetup:   func(_ *mockSvc) {},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "duplicate email",
			body: domain.CreateUserInput{Name: "Bob", Email: "taken@example.com", Role: domain.RoleUser},
			svcSetup: func(s *mockSvc) {
				s.createFn = func(_ context.Context, _ domain.CreateUserInput) (*domain.User, error) {
					return nil, domain.ErrEmailAlreadyTaken
				}
			},
			wantStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			svc := &mockSvc{}
			tt.svcSetup(svc)
			h := newHandler(svc)

			var body *bytes.Buffer
			if s, ok := tt.body.(string); ok {
				body = bytes.NewBufferString(s)
			} else {
				body = toJSON(t, tt.body)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", body)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.Create(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d (body: %s)", w.Code, tt.wantStatus, w.Body.String())
			}
		})
	}
}

func TestDelete_Handler(t *testing.T) {
	t.Parallel()

	svc := &mockSvc{
		deleteFn: func(_ context.Context, id int64) error {
			if id == 999 {
				return domain.ErrUserNotFound
			}
			return nil
		},
		getByIDFn: func(_ context.Context, id int64) (*domain.User, error) {
			if id == 999 {
				return nil, domain.ErrUserNotFound
			}
			return &domain.User{ID: id}, nil
		},
	}

	h := newHandler(svc)
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /api/v1/users/{id}", h.Delete)

	tests := []struct {
		path       string
		wantStatus int
	}{
		{"/api/v1/users/1", http.StatusNoContent},
		{"/api/v1/users/999", http.StatusNotFound},
		{"/api/v1/users/abc", http.StatusBadRequest},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.path, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodDelete, tt.path, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if w.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
