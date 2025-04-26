package order

import (
	"4-order-api/internal/product"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(order *Order) error {
	return r.db.Create(order).Error
}

func (r *Repository) GetByID(id uint64) (*Order, error) {
	var order Order
	if err := r.db.Preload("Products").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *Repository) Update(order *Order) error {
	return r.db.Save(order).Error
}

func (r *Repository) GetByUserID(userID uint64) ([]Order, error) {
	var orders []Order
	if err := r.db.Preload("Products").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *Repository) GetProduct(id uint) (*product.Product, error) {
	var prod product.Product
	if err := r.db.First(&prod, id).Error; err != nil {
		return nil, err
	}
	return &prod, nil
}

func (r *Repository) UpdateProduct(prod *product.Product) error {
	return r.db.Save(prod).Error
}

func (r *Repository) CreateOrderProduct(op *OrderProduct) error {
	return r.db.Create(op).Error
}