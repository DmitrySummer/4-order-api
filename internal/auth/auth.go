package auth

import (
	"4-order-api/config"
	"4-order-api/pkg/db"
	"4-order-api/pkg/jwt"
	"4-order-api/pkg/request"
	resp "4-order-api/pkg/res"
	"net/http"
)

type AuthHandlerDeps struct {
	Config *config.Config
	DB     *db.Db
	Auth   *AuthService
}

type AuthHandler struct {
	config *config.Config
	db     *db.Db
	auth   *AuthService
}

func NewAuthHandler(mux *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		config: deps.Config,
		db:     deps.DB,
		auth:   deps.Auth,
	}

	mux.HandleFunc("/api/auth/login", handler.handleLogin)
	mux.HandleFunc("/api/auth/register", handler.handleRegister)
}

func (h *AuthHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Не допустимый метод", http.StatusMethodNotAllowed)
		return
	}
	body, err := request.HandleBody[LoginRequest](&w, r)
	if err != nil {
		return
	}
	phone, err := h.auth.Login(body.Phone, body.Password)
	token, err := jwt.NewJwt(h.config.Auth.Secret).Create(jwt.JWTData{
		Phone: phone,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data := LoginResponse{
		Token: token,
	}
	resp.Json(w, data, 200)

}

func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Не допустимый метод", http.StatusMethodNotAllowed)
		return
	}
	body, err := request.HandleBody[RegisterRequest](&w, r)
	if err != nil {
		return
	}
	phone, err := h.auth.Register(body.Phone, body.Name, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
	token, err := jwt.NewJwt(h.config.Auth.Secret).Create(jwt.JWTData{
		Phone: phone,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data := RegisterResponse{
		Token: token,
	}
	resp.Json(w, data, 200)

}
