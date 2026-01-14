package storage

import ("learn-go/apidemo/internal/utils")

type Storage interface {

	CreateStudent(name string, age int, email string) (int, error)
	GetStudentByID(id int) (utils.Student, error)
}