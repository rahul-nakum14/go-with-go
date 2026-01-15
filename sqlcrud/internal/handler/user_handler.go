package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/rahulnakum/sqlcrud/internal/model"
	"github.com/rahulnakum/sqlcrud/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var u model.User
	json.NewDecoder(r.Body).Decode(&u)

	h.service.Create(&u)
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	u, _ := h.service.Get(id)
	json.NewEncoder(w).Encode(u)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	var u model.User
	json.NewDecoder(r.Body).Decode(&u)
	h.service.Update(&u)
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	h.service.Delete(id)
	w.WriteHeader(http.StatusNoContent)
}
