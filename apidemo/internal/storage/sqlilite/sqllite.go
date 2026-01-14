package sqlilite

import (
	"database/sql"
	"learn-go/apidemo/internal/config"
	_ "github.com/mattn/go-sqlite3"
	"log"
"learn-go/apidemo/internal/utils"
	"os"
)


type Sqlilite struct {
	DB *sql.DB
}

func New(cfg *config.Config) (*Sqlilite, error) {
		logger := log.New(os.Stdout, "http: ", log.LstdFlags)

	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	res, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL,
		email TEXT NOT NULL UNIQUE
	);
	`)

	logger.Println("Table creation result: ", res)
	if err != nil {
		return nil, err
	}
	return &Sqlilite{DB: db}, nil
}

func (s *Sqlilite) CreateStudent(name string, age int, email string) (int, error) {
	stat, err := s.DB.Prepare("INSERT INTO students (name, age, email) VALUES (?, ?, ?)")

	if err != nil {
		return 0, err
	}

	result, err := stat.Exec(name, age, email)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
	// result, err := s.DB.Exec(
	// 	"INSERT INTO students (name, age, email) VALUES (?, ?, ?)",
	// 	name,
	// 	age,
	// 	email,
	// )
	// if err != nil {
	// 	return 0, err
	// }
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }
	// return int(id), nil
	// return 0, nil
}

func (s *Sqlilite) GetStudentByID(id int) (utils.Student, error) {
	var student utils.Student

	err := s.DB.QueryRow(
		"SELECT id, name, age, email FROM students WHERE id = ?",
		id,
	).Scan(
		&student.ID,
		&student.Name,
		&student.Age,
		&student.Email,
	)

	if err != nil {
		return utils.Student{}, err
	}

	return student, nil
}
