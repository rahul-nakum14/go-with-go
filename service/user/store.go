package user

import "database/sql"

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	var user types.User
	err := s.db.QueryRow(
		"SELECT *FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) GetUserByID(id int) (*types.User, error) {
	// var user types.User
	// err := s.db.QueryRow(
	// 	"SELECT *FROM users WHERE id = ?",
	// 	id,
	// ).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt)
	// if err != nil {
	// 	return nil, err
	// }
	// return &user, nil
	return nil, nil
}

func (s *Store) CreateUser(user types.User) error {
	return nil
}