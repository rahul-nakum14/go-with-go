package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"learn-go/apidemo/internal/storage"
	"learn-go/apidemo/internal/utils"
	"net/http"
	"github.com/go-playground/validator/v10"
    "strconv"
)

func Create(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student utils.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			if errors.Is(err, io.EOF) {
				utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Errorf("empty request body"))
				return
			}
			utils.SendResponse(w, http.StatusBadRequest, "Invalid input data", nil, err.Error())
			return
		}

		if err := student.Validate(); err != nil {
			if verrs, ok := err.(validator.ValidationErrors); ok {
				resp := utils.ValidateErrorResponse(verrs)
				utils.SendErrorResponse(w, http.StatusBadRequest, fmt.Errorf(resp.Error))
				return // important!
			}

			utils.SendErrorResponse(w, http.StatusBadRequest, err)
			return
		}

        id, err := storage.CreateStudent(student.Name, student.Age, student.Email)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err)
			return
		}

		utils.SendResponse(w, http.StatusCreated, "Student created successfully", map[string]interface{}{"id": id}, "")    
	}
}

func GetByID(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		
		studentID, err := strconv.Atoi(id)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		student, err := storage.GetStudentByID(studentID)
		if err != nil {
			utils.SendErrorResponse(w, http.StatusInternalServerError, err)
			return
		}
		
		utils.SendResponse(w, http.StatusOK, "Student retrieved successfully", student, "")
	}
}