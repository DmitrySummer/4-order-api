package order

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"4-order-api/internal/product"
	"4-order-api/pkg/middleware"
	"4-order-api/pkg/res"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux, auth *middleware.AuthMiddleware) {
	mux.Handle("/orders", auth.RequireAuth(http.HandlerFunc(h.CreateOrder)))
	mux.Handle("/orders/", auth.RequireAuth(http.HandlerFunc(h.GetOrder)))
	mux.Handle("/orders/my", auth.RequireAuth(http.HandlerFunc(h.GetMyOrders)))
}

func (h *Handler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		res.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var input struct {
		Products []struct {
			ID       uint `json:"id" validate:"required"`
			Quantity int  `json:"quantity" validate:"required,gt=0"`
		} `json:"products" validate:"required,dive"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		res.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		res.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	tx := h.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	order := Order{
		UserID: userID,
		Status: "pending",
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		res.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	var totalPrice float64
	for _, p := range input.Products {
		var prod product.Product
		if err := tx.First(&prod, p.ID).Error; err != nil {
			tx.Rollback()
			res.Error(w, http.StatusNotFound, err.Error())
			return
		}

		if prod.Stock < p.Quantity {
			tx.Rollback()
			res.Error(w, http.StatusBadRequest, "insufficient stock")
			return
		}

		prod.Stock -= p.Quantity
		if err := tx.Save(&prod).Error; err != nil {
			tx.Rollback()
			res.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		if err := tx.Create(&OrderProduct{
			OrderID:   order.ID,
			ProductID: prod.ID,
			Quantity:  p.Quantity,
		}).Error; err != nil {
			tx.Rollback()
			res.Error(w, http.StatusInternalServerError, err.Error())
			return
		}

		totalPrice += float64(p.Quantity) * prod.Price
	}

	order.TotalPrice = totalPrice
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		res.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		res.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	res.Success(w, order)
}

func (h *Handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		res.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/orders/")
	orderID, err := strconv.ParseUint(path, 10, 64)
	if err != nil {
		res.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		res.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var order Order
	if err := h.db.Preload("Products").First(&order, orderID).Error; err != nil {
		res.Error(w, http.StatusNotFound, err.Error())
		return
	}

	if order.UserID != userID {
		res.Error(w, http.StatusForbidden, "access denied")
		return
	}

	res.Success(w, order)
}

func (h *Handler) GetMyOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		res.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID := middleware.GetUserIDFromContext(r.Context())
	if userID == 0 {
		res.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var orders []Order
	if err := h.db.Where("user_id = ?", userID).Preload("Products").Find(&orders).Error; err != nil {
		res.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	res.Success(w, orders)
}
