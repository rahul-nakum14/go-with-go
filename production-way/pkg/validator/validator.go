// Package validator wraps go-playground/validator with a singleton instance.
package validator

import (
	"sync"

	"github.com/go-playground/validator/v10"
)

// Validator is a thin wrapper around go-playground/validator.
type Validator struct {
	v *validator.Validate
}

var (
	instance *Validator
	once     sync.Once
)

// New returns the shared Validator instance (initialised once).
func New() *Validator {
	once.Do(func() {
		v := validator.New(validator.WithRequiredStructFields())
		// Register custom validations here if needed.
		instance = &Validator{v: v}
	})
	return instance
}

// Struct validates a struct and returns any validation errors.
func (vl *Validator) Struct(s any) error {
	return vl.v.Struct(s)
}
