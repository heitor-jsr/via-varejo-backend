package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"via-varejo/internal/domain"
	"via-varejo/routes"

	_ "github.com/jackc/pgx"
)

const port = "8080"

func main() {
	pool, err := domain.ConnectToDB()
	if err != nil {
		log.Fatalf("Não foi possível conectar ao banco de dados: %v\n", err)
	}
	defer pool.Close()

	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Falha ao pingar o banco de dados: %v\n", err)
	}

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
