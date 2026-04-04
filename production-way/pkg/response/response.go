// Package response provides helpers for writing consistent JSON HTTP responses.
package response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// Error is the standard error envelope.
type Error struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// PaginatedResponse is the standard paginated list envelope.
type PaginatedResponse struct {
	Data       any   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	TotalPages int64 `json:"total_pages"`
}

// JSON writes v as JSON with the given status code.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Last resort — headers are already written, nothing more we can do.
		http.Error(w, "encoding error", http.StatusInternalServerError)
	}
}

// Error writes a JSON error response.
func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, &Error{Code: status, Message: message})
}

// InternalError writes a generic 500 without leaking internals.
func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "an unexpected error occurred")
}

// ValidationError unpacks validator.ValidationErrors into a field→message map.
func ValidationError(w http.ResponseWriter, err error) {
	details := make(map[string]string)
	var ve validator.ValidationErrors
	if ok := errorAs(err, &ve); ok {
		for _, fe := range ve {
			details[fieldName(fe.Field())] = validationMessage(fe)
		}
	}
	JSON(w, http.StatusUnprocessableEntity, &Error{
		Code:    http.StatusUnprocessableEntity,
		Message: "validation failed",
		Details: details,
	})
}

// Paginated writes a paginated list response.
func Paginated(w http.ResponseWriter, data any, total int64, page, perPage int) {
	pages := total / int64(perPage)
	if total%int64(perPage) != 0 {
		pages++
	}
	JSON(w, http.StatusOK, &PaginatedResponse{
		Data:       data,
		Total:      total,
		Page:       page,
		PerPage:    perPage,
		TotalPages: pages,
	})
}

// Decode reads and decodes JSON from the request body.
func Decode(r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(nil, r.Body, 1<<20) // 1 MB limit
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(v); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

// ── Internal helpers ──────────────────────────────────────────────────────────

func errorAs[T error](err error, target *T) bool {
	// We do a type assertion rather than errors.As to keep the import minimal.
	if t, ok := err.(T); ok {
		*target = t
		return true
	}
	return false
}

func fieldName(f string) string {
	// lowercase first char: "Name" → "name"
	if len(f) == 0 {
		return f
	}
	return string(f[0]+32) + f[1:]
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", fe.Param())
	default:
		return fe.Error()
	}
}
