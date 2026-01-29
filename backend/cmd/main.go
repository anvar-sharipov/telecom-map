package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/anvar-sharipov/telecom-map/internal/db"
	"github.com/anvar-sharipov/telecom-map/internal/handler"

	// "github.com/anvar-sharipov/telecom-map/internal/handler/admin"

	admingroups "github.com/anvar-sharipov/telecom-map/internal/handler/admin/groups"
	"github.com/anvar-sharipov/telecom-map/internal/handler/admin/users"
	"github.com/anvar-sharipov/telecom-map/internal/middleware"
	"github.com/anvar-sharipov/telecom-map/internal/repository"
	"github.com/anvar-sharipov/telecom-map/internal/repository/postgres"
	"github.com/anvar-sharipov/telecom-map/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	// w prod ne nujen wstawlyaet dannye w os.Getenv w prode w os.Getenv wstawlyaet dannye pri docker-compose w ney est env_file punkt tam wybiratesya env
	appEnv := os.Getenv("APP_ENV")

	if appEnv != "prod" {
		err := godotenv.Load("../.env.local")
		if err != nil {
			log.Fatal("Error loading ../.env.local file")
		}
	}

	// err := godotenv.Load(".env.local")
	// if err != nil {
	// 	log.Fatal("Error loading .env.local file")
	// }

	APP_PORT := os.Getenv("APP_PORT")
	if APP_PORT == "" {
		APP_PORT = "8000"
	}

	apiPrefix := os.Getenv("API_PREFIX")

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		if appEnv == "prod" {
			frontendURL = "https://192.168.25.74"
		} else {
			frontendURL = "http://localhost:5173"
		}
	}

	pool, err := db.NewPostgresPool()
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer pool.Close()

	fmt.Println("✅ Connected to Postgres successfully!")

	// 1️⃣ СНАЧАЛА создаём репозитории
	userRepo := postgres.NewUserRepository(pool)
	groupRepo := repository.NewGroupRepository(pool)
	refreshTokenRepo := repository.NewRefreshTokenRepository(pool) // 🔥 ТУТ БЫЛО refreshRepo

	// 2️⃣ ПОТОМ создаём сервис (использует refreshTokenRepo)
	authService := &service.AuthService{
		RefreshTokenRepo: refreshTokenRepo, // 🔥 ТЕПЕРЬ ОН УЖЕ СУЩЕСТВУЕТ
	}

	// authHandler := &handler.AuthHandler{UserRepo: userRepo}
	// refreshRepo := repository.NewRefreshTokenRepository(pool)
	authHandler := &handler.AuthHandler{
		UserRepo:         userRepo,
		RefreshTokenRepo: refreshTokenRepo,
		AuthService:      authService,
	}

	adminUsersHandler := &users.Handler{
		UserRepo: userRepo,
	}

	adminGroupHandler := &admingroups.AdminGroupHandler{
		GroupRepo: groupRepo,
	}

	mux := http.NewServeMux()
	mux.HandleFunc(apiPrefix+"/register", middleware.ErrorMiddleware(authHandler.Register))
	mux.HandleFunc(apiPrefix+"/login", middleware.ErrorMiddleware(authHandler.Login))
	mux.HandleFunc(apiPrefix+"/auth/refresh", middleware.ErrorMiddleware(authHandler.Refresh))
	mux.HandleFunc(apiPrefix+"/auth/logout", middleware.ErrorMiddleware(authHandler.Logout))
	mux.HandleFunc(apiPrefix+"/auth/me", middleware.ErrorMiddleware(authHandler.Me))
	// admin
	mux.HandleFunc(apiPrefix+"/admin/users", middleware.ErrorMiddleware(middleware.Auth(adminUsersHandler.Handle)))
	mux.HandleFunc(apiPrefix+"/admin/groups", middleware.ErrorMiddleware(middleware.Auth(adminGroupHandler.ListGroups)))
	// mux.HandleFunc(apiPrefix+"/admin/users", middleware.ErrorMiddleware(adminHandler.ListUsers))

	// Middleware для CORS
	handlerWithCORS := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin == frontendURL {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			// origin := r.Header.Get("Origin")
			// if appEnv == "prod" && origin == "https://192.168.25.74" {
			// 	w.Header().Set("Access-Control-Allow-Origin", origin)
			// } else if appEnv != "prod" && origin == "http://localhost:5173" {
			// 	w.Header().Set("Access-Control-Allow-Origin", origin)
			// }

			// if appEnv == "prod" {
			// 	w.Header().Set("Access-Control-Allow-Origin", "https://192.168.25.74")
			// } else {
			// 	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
			// }

			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true") // 🔥 ВАЖНО dlya cookie

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}

	fmt.Printf("Server running on :%s\n", APP_PORT)
	log.Printf("🚀 Server started on :%s (env=%s)", APP_PORT, os.Getenv("APP_ENV"))

	log.Fatal(http.ListenAndServe(":"+APP_PORT, handlerWithCORS(mux)))
}
