// domain/db.go

package domain

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

var db *pgxpool.Pool

const dbTimeout = time.Second * 3

func InitDB() error {
	var err error
	db, err = ConnectToDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	return nil
}

func ConnectToDB() (*pgxpool.Pool, error) {
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
		return nil, fmt.Errorf("failed to create a config: %v", err)
	}

	dbConfig.MaxConns = defaultMaxConns
	dbConfig.MinConns = defaultMinConns
	dbConfig.MaxConnLifetime = defaultMaxConnLifetime
	dbConfig.MaxConnIdleTime = defaultMaxConnIdleTime
	dbConfig.HealthCheckPeriod = defaultHealthCheckPeriod
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	pool, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	return pool, nil
}

func GetDB() *pgxpool.Pool {
	return db
}
