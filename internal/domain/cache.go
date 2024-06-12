package domain

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var C *redis.Client

func InitRedisClient(addr, password string, db int) {
	C = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	_, err := C.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Falha ao pingar o servidor Redis: %v\n", err)
	}

	log.Println("Conexão com o servidor Redis estabelecida com sucesso!")
}
