package main

import (
	"log"
	"net/http"
	"time"

	"setho-pymnt/internal/db"
	"setho-pymnt/internal/handlers"
	"setho-pymnt/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	store, err := db.New("data/payment.db")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = store.Close() }()

	if err := store.Migrate(); err != nil {
		log.Fatal(err)
	}

	serviceLayer := service.New(store)
	handler := handlers.New(serviceLayer)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/businesses", handler.ListBusinesses)
		r.Post("/businesses", handler.CreateBusiness)
		r.Get("/businesses/{id}", handler.GetBusiness)

		r.Get("/transactions", handler.ListTransactions)
		r.Post("/transactions", handler.CreateTransaction)
		r.Get("/transactions/{id}", handler.GetTransaction)

		r.Get("/summary", handler.Summary)
	})

	log.Println("Payment system running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
