package main

import (
	"fmt"
	"log"
	"net/http"
	"via-varejo/internal/domain"
	"via-varejo/routes"

	_ "github.com/jackc/pgx"
)

const port = "8080"

func main() {
	err := domain.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer domain.GetDB().Close()

	domain.InitRedisClient("redis:6379", "", 0)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: routes.Routers(),
	}

	fmt.Println("Starting server on port " + port)

	err = srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
