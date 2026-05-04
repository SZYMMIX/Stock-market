package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/szymmix/stock-market/internal/handlers"
	"github.com/szymmix/stock-market/internal/storage"
)

func main() {
	port := os.Getenv("APP_PORT")
	dbURL := os.Getenv("DB_URL")

	if port == "" {
		port = "8080"
	}

	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	db, err := storage.NewPostgresStorage(ctx, dbURL)
	if err != nil {
		log.Fatalf("Database initialization error: %v", err)
	}
	defer db.Close()

	stocksHandler := handlers.NewStocksHandler(db)
	walletsHandler := handlers.NewWalletsHandler(db)
	logsHandler := &handlers.LogsHandler{Storage: db}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/stocks", stocksHandler.SetStocks)
	r.Get("/stocks", stocksHandler.GetStocks)

	r.Get("/wallets/{wallet_id}", walletsHandler.GetWallet)
	r.Get("/wallets/{wallet_id}/stocks/{stock_name}", walletsHandler.GetStockQuantity)
	r.Post("/wallets/{wallet_id}/stocks/{stock_name}", walletsHandler.TradeStock)

	r.Get("/log", logsHandler.GetLogs)

	r.Post("/chaos", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Chaos: Killing an instance...")
		os.Exit(1)
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		fmt.Printf("Stock server starting on port %s...\n", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-done
	fmt.Println("\nClosing server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Error while shutting down the server: %v", err)
	}

	fmt.Println("Server safely shut down.")
}
