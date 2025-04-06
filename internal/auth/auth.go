package auth

import (
	"4-order-api/config"
	"4-order-api/pkg/db"
	"net/http"
)

type AuthHandlerDeps struct {
	Config *config.Config
	DB     *db.Db
}

type AuthHandler struct {
	config *config.Config
	db     *db.Db
}

func NewAuthHandler(mux *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		config: deps.Config,
		db:     deps.DB,
	}

	mux.HandleFunc("/api/auth/login", handler.handleLogin)
	mux.HandleFunc("/api/auth/register", handler.handleRegister)
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Не допустимый метод", http.StatusMethodNotAllowed)
		return
	}

}

func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Не допустимый метод", http.StatusMethodNotAllowed)
		return
	}
}
