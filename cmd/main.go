package main

import (
	"log"
	"shop/internal/delivery/http"
	"shop/internal/repository/mysql"
	"shop/internal/usecase"
	"shop/pkg/database"
)

func main() {
	// Initialize database
	db, err := database.NewMySQLConnection()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := mysql.NewUserRepository(db)
	productRepo := mysql.NewProductRepository(db)
	orderRepo := mysql.NewOrderRepository(db)
	cartRepo := mysql.NewCartRepository(db)
	paymentRepo := mysql.NewPaymentRepository(db)

	// Initialize use cases
	userUsecase := usecase.NewUserUsecase(userRepo)
	productUsecase := usecase.NewProductUsecase(productRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, productRepo)
	cartUsecase := usecase.NewCartUseCase(cartRepo, userRepo, orderRepo, productRepo)
	// TODO Add other usecase to payment for logic
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo, orderRepo, productRepo, userRepo, cartRepo)
	// Initialize HTTP handler
	handler := http.NewHandler(userUsecase, productUsecase, orderUsecase, cartUsecase, paymentUsecase)

	// Start server
	log.Println("Server starting on :8080")
	if err := handler.Router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
