package repository

import (
	"github.com/rahulnakum/sqlcrud/internal/model"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByID(id int64) (*model.User, error)
	Update(user *model.User) error
	Delete(id int64) error
}
