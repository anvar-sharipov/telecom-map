package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anvar-sharipov/telecom-map/internal/db"
	"github.com/anvar-sharipov/telecom-map/internal/handler"
	"github.com/anvar-sharipov/telecom-map/internal/middleware"
	"github.com/anvar-sharipov/telecom-map/internal/repository/postgres"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env.local file")
	}

	APP_PORT := os.Getenv("APP_PORT")

	pool, err := db.NewPostgresPool()
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer pool.Close()

	fmt.Println("✅ Connected to Postgres successfully!")

	userRepo := postgres.NewUserRepository(pool)
	authHandler := &handler.AuthHandler{UserRepo: userRepo}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", middleware.ErrorMiddleware(authHandler.Register))
	mux.HandleFunc("/login", middleware.ErrorMiddleware(authHandler.Login))

	// Middleware для CORS
	handlerWithCORS := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}

	fmt.Printf("Server running on :%s\n", APP_PORT)
	log.Fatal(http.ListenAndServe(":"+APP_PORT, handlerWithCORS(mux)))
}
