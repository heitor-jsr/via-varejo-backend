package routes

import (
	"fmt"
	"net/http"
	"via-varejo/internal/app"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Routers() http.Handler {
	fmt.Println("routers is running")
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Use(middleware.Heartbeat("/ping"))
	mux.Use(middleware.Logger)

	mux.Post("/purchase", app.CreateNewPurchaseSummary)

	mux.Get("/purchase/retrivedPurchaseSummary", app.GetRedisPurchaseSummary)
	return mux
}
