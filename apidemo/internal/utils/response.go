package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	statusOk 	 = "OK"
	statusError  = "ERROR"
)

func SendResponse(w http.ResponseWriter, statusCode int, message string, data interface{}, err string) {
	response := APIResponse{
		Message: message,
		Data:    data,
		Error:   err,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)


	if jsonErr := json.NewEncoder(w).Encode(response); jsonErr != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func SendErrorResponse(
	w http.ResponseWriter,
	statusCode int,
	err error,
) {
	SendResponse(
		w,
		statusCode,
		statusError,
		nil,
		err.Error(),
	)
}

func ValidateErrorResponse(
	err validator.ValidationErrors,
) APIResponse {
	var errorMessages []string
	for _, validationErr := range err {
		switch validationErr.ActualTag() {
		case "required":
			errorMessages = append(
				errorMessages,
				fmt.Sprintf("The field '%s' is required.", validationErr.Field()),	
			)
		
		case "email":
			errorMessages = append(
				errorMessages,
				fmt.Sprintf("The field '%s' must be a valid email address.", validationErr.Field()),	
			)
		case "gte":
			errorMessages = append(
				errorMessages,
				fmt.Sprintf("The field '%s' must be greater than or equal to %s.", validationErr.Field(), validationErr.Param()),	
			)
		case "lte":
			errorMessages = append(
				errorMessages,
				fmt.Sprintf("The field '%s' must be less than or equal to %s.", validationErr.Field(), validationErr.Param()),	
			)
		}

	}
	return APIResponse{
		Message: statusError,
		Error:   strings.Join(errorMessages, " "),
	}
}