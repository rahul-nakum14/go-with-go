package student

import (
	"encoding/json"
	"errors"
	"io"
	"learn-go/apidemo/internal/utils"
	"net/http"
)
 
func Create(w http.ResponseWriter, r *http.Request) {
    var student utils.Student

    if r.Body == nil {
        utils.SendResponse(w, http.StatusBadRequest, "Request body is empty", nil, "Body is empty")
        return
    }

    err := json.NewDecoder(r.Body).Decode(&student)
    if err != nil {
        if errors.Is(err, io.EOF) {
            utils.SendResponse(w, http.StatusBadRequest, "Request body is empty", nil, "Body is empty")
            return
        }
        utils.SendResponse(w, http.StatusBadRequest, "Invalid input data", nil, err.Error())
        return
    }

    if err := student.Validate(); err != nil {
        utils.SendResponse(w, http.StatusBadRequest, "Invalid input data", nil, err.Error())
        return
    }

    utils.SendResponse(w, http.StatusCreated, "Student created successfully", student, "")
}