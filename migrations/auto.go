package main

import (
	"fmt"
	"os"

	"4-order-api/internal/product"
	"4-order-api/internal/user"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&product.Product{}, user.User{})
	if err != nil {
		fmt.Printf("Ошибка при выполнении миграции: %v\n", err)
		panic(err)
	}

	fmt.Println("Миграция успешно выполнена")
}
