package repository

import (
	"database/sql"
	"github.com/rahulnakum/sqlcrud/internal/model"
)

type userRepo struct {
	db     *sql.DB
	driver string
}

func NewUserRepo(db *sql.DB, driver string) UserRepository {
	return &userRepo{db: db, driver: driver}
}

func (r *userRepo) Create(u *model.User) error {
	if r.driver == "postgres" {
		return r.db.QueryRow(
			`INSERT INTO users (name, email) VALUES ($1,$2) RETURNING id`,
			u.Name, u.Email,
		).Scan(&u.ID)
	}

	res, err := r.db.Exec(
		`INSERT INTO users (name, email) VALUES (?,?)`,
		u.Name, u.Email,
	)
	if err != nil {
		return err
	}
	u.ID, _ = res.LastInsertId()
	return nil
}

func (r *userRepo) GetByID(id int64) (*model.User, error) {
	u := &model.User{}
	query := `SELECT id,name,email FROM users WHERE id=?`
	if r.driver == "postgres" {
		query = `SELECT id,name,email FROM users WHERE id=$1`
	}
	err := r.db.QueryRow(query, id).Scan(&u.ID, &u.Name, &u.Email)
	return u, err
}

func (r *userRepo) Update(u *model.User) error {
	query := `UPDATE users SET name=?, email=? WHERE id=?`
	if r.driver == "postgres" {
		query = `UPDATE users SET name=$1, email=$2 WHERE id=$3`
	}
	_, err := r.db.Exec(query, u.Name, u.Email, u.ID)
	return err
}

func (r *userRepo) Delete(id int64) error {
	query := `DELETE FROM users WHERE id=?`
	if r.driver == "postgres" {
		query = `DELETE FROM users WHERE id=$1`
	}
	_, err := r.db.Exec(query, id)
	return err
}
