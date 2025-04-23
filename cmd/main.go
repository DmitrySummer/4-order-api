package main

import (
	"4-order-api/config"
	"4-order-api/internal/auth"
	"4-order-api/internal/handler"
	"4-order-api/internal/order"
	"4-order-api/internal/user"
	"4-order-api/pkg/db"
	"4-order-api/pkg/middleware"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	conf := config.LoadConfig()

	database, err := db.NewDb(conf)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	router := http.NewServeMux()

	userRepository := user.NewRepository(database.DB)
	authService := auth.NewAuthService(userRepository)
	authMiddleware := middleware.NewAuthMiddleware(conf, userRepository)

	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
		DB:     database,
		Auth:   authService,
	})

	handler.NewProductHandler(router, handler.ProductHandlerDeps{
		Config:         conf,
		DB:            database,
		UserRepository: userRepository,
	})

	handler.NewUserHandler(router, handler.UserHandlerDeps{
		DB:          database,
		AuthService: authService,
	})

	orderHandler := order.NewHandler(database.DB)
	orderHandler.RegisterRoutes(router, authMiddleware)

	server := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("Server listening on port 8081")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	fmt.Println("\nShutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	fmt.Println("Server gracefully stopped")
}
