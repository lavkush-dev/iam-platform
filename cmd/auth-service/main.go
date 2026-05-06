package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"iam-platform/internal/logger"
	appmw "iam-platform/internal/middleware"

	"iam-platform/internal/config"
	"iam-platform/internal/handlers"
	"iam-platform/internal/repository/postgres"
	"iam-platform/internal/services"
	"iam-platform/pkg/jwt"
)

func main() {
	cfg := config.Load()

	log := logger.New()
	defer log.Sync()

	db := postgres.NewDB(cfg.DB.URL)

	userRepo := postgres.NewUserRepository(db)

	jwtManager := jwt.NewManager(cfg.JWT.Secret)

	authService := services.NewAuthService(userRepo, jwtManager)

	authHandler := handlers.NewAuthHandler(authService, log)

	r := chi.NewRouter()

	// middleware
	r.Use(appmw.ZapLogger(log))

	r.Post("/login", authHandler.Login)

	log.Info("auth-service running") // optional structured field style

	http.ListenAndServe(":"+cfg.Server.Port, r)
}
