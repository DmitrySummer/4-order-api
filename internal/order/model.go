package order

import (
	"4-order-api/internal/product"
	"4-order-api/internal/user"
	"time"
)

type OrderProduct struct {
	OrderID   uint `gorm:"primaryKey"`
	ProductID uint `gorm:"primaryKey"`
	Quantity  int  `json:"quantity" validate:"required,gt=0"`
}

type Order struct {
	ID         uint              `json:"id" gorm:"primaryKey"`
	UserID     uint              `json:"user_id" validate:"required"`
	User       user.User         `json:"user" gorm:"foreignKey:UserID"`
	Products   []product.Product `json:"products" gorm:"many2many:order_products;"`
	TotalPrice float64           `json:"total_price"`
	Status     string            `json:"status" validate:"required,oneof=pending processing completed cancelled"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
}
