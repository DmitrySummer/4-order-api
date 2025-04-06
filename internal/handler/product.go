package handler

import (
	"4-order-api/internal/product"
	"4-order-api/pkg/db"
	"encoding/json"
	"net/http"
)

type ProductHandlerDeps struct {
	DB *db.Db
}

type ProductHandler struct {
	db *db.Db
}

func NewProductHandler(mux *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		db: deps.DB,
	}

	mux.HandleFunc("/api/products", handler.handleProducts)
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
