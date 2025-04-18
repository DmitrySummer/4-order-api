package handler

import (
	"4-order-api/config"
	"4-order-api/internal/product"
	"4-order-api/internal/user"
	"4-order-api/pkg/db"
	"4-order-api/pkg/middleware"
	"encoding/json"
	"net/http"
	"strconv"
)

type ProductHandlerDeps struct {
	Config         *config.Config
	DB             *db.Db
	UserRepository *user.Repository
}

type ProductHandler struct {
	config         *config.Config
	db             *db.Db
	userRepository *user.Repository
}

func NewProductHandler(mux *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		config:         deps.Config,
		db:             deps.DB,
		userRepository: deps.UserRepository,
	}

	authMiddleware := middleware.NewAuthMiddleware(deps.Config, deps.UserRepository)

	mux.HandleFunc("/api/products", handler.handleProducts)

	mux.Handle("/api/products/buy", authMiddleware.RequireAuth(http.HandlerFunc(handler.handleBuyProduct)))
}

func (h *ProductHandler) handleProducts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getProducts(w, r)
	case http.MethodPost:
		h.createProduct(w, r)
	default:
		http.Error(w, "Не допустимый метод", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) handleBuyProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	currentUser := middleware.GetUserFromContext(r.Context())
	if currentUser == nil {
		http.Error(w, "Пользователь не авторизован", http.StatusUnauthorized)
		return
	}

	productIDStr := r.URL.Query().Get("product_id")
	if productIDStr == "" {
		http.Error(w, "ID продукта не указан", http.StatusBadRequest)
		return
	}

	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Неверный формат ID продукта", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"message":    "Товар успешно куплен",
		"product_id": productID,
		"user_phone": currentUser.Phone,
		"user_name":  currentUser.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ProductHandler) getProducts(w http.ResponseWriter, r *http.Request) {
	products := []product.Product{}
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct product.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}
