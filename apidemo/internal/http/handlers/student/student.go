package student

import (
	"encoding/json"
	"errors"
	"io"
	"learn-go/apidemo/internal/utils"
	"net/http"
    "github.com/go-playground/validator/v10"
    "fmt"
)
 
func Create(w http.ResponseWriter, r *http.Request) {
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
        return  // important!
    }

    utils.SendErrorResponse(w, http.StatusBadRequest, err)
    return
}

}

func Ping(w http.ResponseWriter, r *http.Request) {
	utils.SendResponse(
		w,
		http.StatusOK,
		"Pong",
		nil,
		"",
	)
}
