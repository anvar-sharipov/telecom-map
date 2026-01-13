package db

import (
	"context" // context → управление временем жизни соединений
	"fmt"
	"os" // os → чтение env
	"time"

	"github.com/jackc/pgx/v5/pgxpool" // pgxpool → пул соединений к Postgres
)

func NewPostgresPool() (*pgxpool.Pool, error) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "localhost" {
		dbHost = "127.0.0.1"
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// 🔍 ОТЛАДКА
	fmt.Printf("🔧 Connecting with:\n")
	fmt.Printf("   User: [%s]\n", dbUser)
	fmt.Printf("   Pass: [%s]\n", dbPass)
	fmt.Printf("   Host: [%s]\n", dbHost)
	fmt.Printf("   Port: [%s]\n", dbPort)
	fmt.Printf("   Name: [%s]\n", dbName)

	sslMode := os.Getenv("DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable" // dev по умолчанию
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		dbUser, dbPass, dbHost, dbPort, dbName, sslMode,
	)

	// Выведите полный DSN (БЕЗ пароля в продакшене!)
	fmt.Printf("🔗 DSN: postgres://%s:***@%s:%s/%s?sslmode=disable\n",
		dbUser, dbHost, dbPort, dbName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	return pool, nil
}
