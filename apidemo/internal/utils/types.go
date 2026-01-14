package utils

import (
	"github.com/go-playground/validator/v10"
)

type Student struct {
	ID    int64 `json:"id"`
	Name  string `json:"name" validate:"required"`
	Age   int    `json:"age" validate:"gte=18,lte=120"`
	Email string `json:"email" validate:"email,required"`
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (s *Student) Validate() error {
	return validate.Struct(s)
}

type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}




// func Ping(w http.ResponseWriter, r *http.Request) {
// 	utils.SendResponse(
// 		w,
// 		http.StatusOK,
// 		"Pong",
// 		nil,
// 		"",
// 	)
// }