package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/andreodev/go-pay/internal/config"
	"github.com/andreodev/go-pay/internal/database"
	"github.com/andreodev/go-pay/internal/http/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	db, err := database.NewPostgresDB(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("GoPay Guard API is running with database"))
	})

	cfg := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	r.Route("/webhooks", func(r chi.Router) {
		r.Use(middleware.APIKeyAuth(cfg.APIKey))

		r.Post("/payments", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("webhook received"))
		})
	})

	addr := fmt.Sprintf(":%s", port)

	log.Println("GoPay Guard running on", addr)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
