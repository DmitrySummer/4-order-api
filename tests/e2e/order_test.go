package e2e

import (
	"4-order-api/config"
	"4-order-api/internal/auth"
	"4-order-api/internal/order"
	"4-order-api/internal/product"
	"4-order-api/internal/user"
	"4-order-api/pkg/db"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestOrderE2E(t *testing.T) {
	cfg := &config.Config{
		Db: config.DbConfig{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			DbName:   "order_api_test",
		},
	}

	database, err := db.NewDb(cfg)
	assert.NoError(t, err)
	defer database.Close()

	database.DB.Exec("DELETE FROM order_products")
	database.DB.Exec("DELETE FROM orders")
	database.DB.Exec("DELETE FROM products")
	database.DB.Exec("DELETE FROM users")

	userRepo := user.NewRepository(database.DB)
	authService := auth.NewAuthService(userRepo)

	testUser := &user.User{
		Phone:    "+79991234567",
		Password: "testpass123",
		Name:     "Test User",
	}

	userID, err := authService.Register(testUser.Phone, testUser.Name, testUser.Password)
	assert.NoError(t, err)

	productRepo := product.NewRepository(database.DB)
	testProducts := []product.Product{
		{
			Name:        "Test Product 1",
			Description: "Description 1",
			Price:       100.0,
			Stock:       10,
		},
		{
			Name:        "Test Product 2",
			Description: "Description 2",
			Price:       200.0,
			Stock:       5,
		},
	}

	for i := range testProducts {
		err = productRepo.Create(&testProducts[i])
		assert.NoError(t, err)
		assert.NotZero(t, testProducts[i].ID)
	}

	orderRepo := order.NewRepository(database.DB)
	testOrder := &order.Order{
		UserID:     userID,
		Products:   testProducts,
		TotalPrice: 500.0,
		Status:     "pending",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	orderProducts := []order.OrderProduct{
		{
			ProductID: testProducts[0].ID,
			Quantity:  2,
		},
		{
			ProductID: testProducts[1].ID,
			Quantity:  1,
		},
	}

	err = orderRepo.Create(testOrder)
	assert.NoError(t, err)
	assert.NotZero(t, testOrder.ID)

	for i := range orderProducts {
		orderProducts[i].OrderID = testOrder.ID
		err = database.DB.Create(&orderProducts[i]).Error
		assert.NoError(t, err)
	}

	savedOrder, err := orderRepo.GetByID(uint64(testOrder.ID))
	assert.NoError(t, err)
	assert.Equal(t, testOrder.UserID, savedOrder.UserID)
	assert.Equal(t, testOrder.TotalPrice, savedOrder.TotalPrice)
	assert.Equal(t, testOrder.Status, savedOrder.Status)
	assert.Len(t, savedOrder.Products, len(testProducts))

	savedOrder.Status = "completed"
	err = orderRepo.Update(savedOrder)
	assert.NoError(t, err)

	updatedOrder, err := orderRepo.GetByID(uint64(testOrder.ID))
	assert.NoError(t, err)
	assert.Equal(t, "completed", updatedOrder.Status)

	database.DB.Exec("DELETE FROM order_products WHERE order_id = ?", testOrder.ID)
	database.DB.Exec("DELETE FROM orders WHERE id = ?", testOrder.ID)
	for _, p := range testProducts {
		database.DB.Exec("DELETE FROM products WHERE id = ?", p.ID)
	}
	database.DB.Exec("DELETE FROM users WHERE id = ?", userID)
}
