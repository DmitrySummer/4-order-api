package auth

import (
	"4-order-api/config"
	"4-order-api/pkg/db"
	"4-order-api/pkg/jwt"
	"4-order-api/pkg/request"
	"4-order-api/pkg/res"
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
		res.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	body, err := request.HandleBody[LoginRequest](&w, r)
	if err != nil {
		http.Error(w, "Невозможно обработать тело запроса: "+err.Error(), http.StatusBadRequest)
		return
	}
	phone, err := h.auth.Login(body.Phone, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := jwt.NewJwt(h.config.Auth.Secret).Create(jwt.JWTData{
		Phone: phone,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := LoginResponse{
		Token: token,
	})
}

func (h *AuthHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	body, err := request.HandleBody[RegisterRequest](&w, r)
	if err != nil {
		http.Error(w, "Невозможно обработать тело запроса: "+err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := h.auth.Register(body.Phone, body.Name, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	token, err := jwt.NewJwt(h.config.Auth.Secret).Create(jwt.JWTData{
		Phone: phone,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := RegisterResponse{
		Token: token,
	}
	resp.Json(w, data, 200)

}
