package handler

import (
	"4-order-api/internal/auth"
	"4-order-api/internal/user"
	"4-order-api/pkg/db"
	"encoding/json"
	"net/http"
)

type UserHandlerDeps struct {
	DB          *db.Db
	AuthService *auth.AuthService
}

type UserHandler struct {
	db          *db.Db
	authService *auth.AuthService
}

func NewUserHandler(mux *http.ServeMux, deps UserHandlerDeps) {
	handler := &UserHandler{
		db:          deps.DB,
		authService: deps.AuthService,
	}

	mux.HandleFunc("/api/user", handler.handleUser)
}

func (h *UserHandler) handleUser(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUser(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		http.Error(w, "Не допустимый метод", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) getUser(w http.ResponseWriter, r *http.Request) {
	user := []user.User{}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var newUser user.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}
