package domain

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

func ConfigPGX() *pgxpool.Config {
	const defaultMaxConns = int32(5)
	const defaultMinConns = int32(0)
	const defaultMaxConnLifetime = time.Hour
	const defaultMaxConnIdleTime = time.Minute * 30
	const defaultHealthCheckPeriod = time.Minute
	const defaultConnectTimeout = time.Second * 5

	DATABASE_URL := os.Getenv("DATABASE_URL")
	if DATABASE_URL == "" {
		DATABASE_URL = "postgres://via_varejo_user:123456@postgres:5432/via_varejo"
	}
	dbConfig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		log.Fatal("Failed to create a config, error: ", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return dbConfig
}

func ConnectToDB() (*pgxpool.Pool, error) {
	var pool *pgxpool.Pool
	var err error

	for i := 0; i < 10; i++ {
		log.Printf("Connecting to Postgres, attempt %d\n", i+1)

		pool, err = pgxpool.ConnectConfig(context.Background(), ConfigPGX())
		if err != nil {
			log.Printf("Failed to connect to Postgres: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		log.Println("Connected to Postgres")
		return pool, nil
	}

	return nil, fmt.Errorf("failed to connect to Postgres after 10 attempts")
}
