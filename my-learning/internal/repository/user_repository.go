package repository

import (
	"database/sql"
	"fmt"

	"github.com/rahul-nakum14/go-user-crud/internal/models"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}
func (r *UserRepository) GetAllUsers() ([]models.User,error) {
	rows, err := r.DB.Query("SELECT id, username, email, password FROM users")

	if err != nil {
		return nil, fmt.Errorf("failed to query users: %v", err)
	}
	defer rows.Close()

	var users []models.User

	for rows.Next(){
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepository) CreateUser(user models.User) (int64, error) {
	var id int64
	err := r.DB.QueryRow("INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id",
		user.Username, user.Email, user.Password).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create user: %v", err)
	}
	return id, nil
}