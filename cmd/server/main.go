package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/mytheresa/go-hiring-challenge/app/catalog"
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/mytheresa/go-hiring-challenge/models"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// signal handling for graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Initialize database connection
	db, close := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	defer close()

	// Initialize handlers
	// Note: prodRepo now follows the ProductRepository interface
	prodRepo := models.NewProductsRepository(db)
	cat := catalog.NewCatalogHandler(prodRepo)

	// Set up routing
	mux := http.NewServeMux()

	// Task 1, 4, 5: Catalog List with Pagination & Filtering
	mux.HandleFunc("GET /catalog", cat.HandleGet)

	// Task 2.1: Product Details (the {code} syntax works in Go 1.22+)
	mux.HandleFunc("GET /catalog/{code}", cat.HandleGetByCode)

	// New Routes (Task 3) -  using the same 'cat' handler instance
	mux.HandleFunc("GET /categories", cat.HandleGetCategories)
	mux.HandleFunc("POST /categories", cat.HandleCreateCategory)

	// Set up the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", os.Getenv("HTTP_PORT")),
		Handler: mux,
	}

	// Start the server
	go func() {
		log.Printf("Starting server on http://%s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %s", err)
		}
	}()

	<-ctx.Done()
	log.Println("Shutting down server...")

	// Shutdown gracefully
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
