// Package handler contains HTTP handlers. Each handler decodes the request,
// calls the service, then encodes the response. No business logic lives here.
package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/yourorg/go-user-crud/internal/domain"
	"github.com/yourorg/go-user-crud/pkg/response"
	"github.com/yourorg/go-user-crud/pkg/validator"
)

// userServicer is the interface the handler depends on.
// Depending on an interface (not the concrete service) makes the handler testable.
type userServicer interface {
	GetByID(ctx context.Context, id int64) (*domain.User, error)
	List(ctx context.Context, filter domain.ListFilter) ([]*domain.User, int64, error)
	Create(ctx context.Context, input domain.CreateUserInput) (*domain.User, error)
	Update(ctx context.Context, id int64, input domain.UpdateUserInput) (*domain.User, error)
	Delete(ctx context.Context, id int64) error
}

// UserHandler groups all user-related HTTP handlers.
type UserHandler struct {
	svc      userServicer
	validate *validator.Validator
	log      *slog.Logger
}

// NewUserHandler constructs a UserHandler.
func NewUserHandler(svc userServicer, v *validator.Validator, log *slog.Logger) *UserHandler {
	return &UserHandler{svc: svc, validate: v, log: log}
}

// HealthCheck godoc
// @Summary     Health check
// @Tags        health
// @Success     200 {object} map[string]string
// @Router      /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// List godoc
// @Summary     List users
// @Tags        users
// @Produce     json
// @Param       page     query int    false "Page number"     default(1)
// @Param       per_page query int    false "Items per page"  default(20)
// @Param       search   query string false "Search by name/email"
// @Success     200 {object} response.Paginated
// @Router      /api/v1/users [get]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := domain.ListFilter{
		Page:    queryIntOr(r, "page", 1),
		PerPage: queryIntOr(r, "per_page", 20),
		Search:  r.URL.Query().Get("search"),
	}

	if err := h.validate.Struct(filter); err != nil {
		response.ValidationError(w, err)
		return
	}

	users, total, err := h.svc.List(r.Context(), filter)
	if err != nil {
		h.log.ErrorContext(r.Context(), "List users failed", "error", err)
		response.InternalError(w)
		return
	}

	response.Paginated(w, users, total, filter.Page, filter.PerPage)
}

// GetByID godoc
// @Summary     Get a user
// @Tags        users
// @Produce     json
// @Param       id path int true "User ID"
// @Success     200 {object} domain.User
// @Failure     404 {object} response.Error
// @Router      /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	user, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}
		h.log.ErrorContext(r.Context(), "GetByID failed", "error", err, "id", id)
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// Create godoc
// @Summary     Create a user
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       body body domain.CreateUserInput true "User input"
// @Success     201 {object} domain.User
// @Failure     422 {object} response.Error
// @Router      /api/v1/users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateUserInput
	if err := response.Decode(r, &input); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validate.Struct(input); err != nil {
		response.ValidationError(w, err)
		return
	}

	user, err := h.svc.Create(r.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrEmailAlreadyTaken) {
			response.Error(w, http.StatusConflict, "email already taken")
			return
		}
		h.log.ErrorContext(r.Context(), "Create user failed", "error", err)
		response.InternalError(w)
		return
	}

	response.JSON(w, http.StatusCreated, user)
}

// Update godoc
// @Summary     Update a user
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id   path int                   true "User ID"
// @Param       body body domain.UpdateUserInput true "Update input"
// @Success     200 {object} domain.User
// @Router      /api/v1/users/{id} [put]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	var input domain.UpdateUserInput
	if err := response.Decode(r, &input); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validate.Struct(input); err != nil {
		response.ValidationError(w, err)
		return
	}

	user, err := h.svc.Update(r.Context(), id, input)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			response.Error(w, http.StatusNotFound, "user not found")
		case errors.Is(err, domain.ErrEmailAlreadyTaken):
			response.Error(w, http.StatusConflict, "email already taken")
		default:
			h.log.ErrorContext(r.Context(), "Update user failed", "error", err, "id", id)
			response.InternalError(w)
		}
		return
	}

	response.JSON(w, http.StatusOK, user)
}

// Delete godoc
// @Summary     Delete a user
// @Tags        users
// @Param       id path int true "User ID"
// @Success     204
// @Failure     404 {object} response.Error
// @Router      /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := pathID(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid user id")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.Error(w, http.StatusNotFound, "user not found")
			return
		}
		h.log.ErrorContext(r.Context(), "Delete user failed", "error", err, "id", id)
		response.InternalError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// pathID extracts the {id} path value (Go 1.22+ stdlib routing).
func pathID(r *http.Request) (int64, error) {
	return strconv.ParseInt(r.PathValue("id"), 10, 64)
}

func queryIntOr(r *http.Request, key string, fallback int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n < 1 {
		return fallback
	}
	return n
}
