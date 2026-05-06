package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"iam-platform/internal/config"
	"iam-platform/internal/handlers"
	"iam-platform/internal/repository/postgres"
	"iam-platform/internal/services"
)

func main() {
	cfg := config.Load()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db := postgres.NewDB(cfg.DB.URL)

	userRepo := postgres.NewUserRepository(db)
	userService := services.NewUserService(userRepo)

	userHandler := handlers.NewUserHandler(userService, logger)

	r := chi.NewRouter()

	r.Post("/users", userHandler.CreateUser)
	r.Get("/users/{id}", userHandler.GetUser)

	logger.Info("iam-service running on :8081")

	http.ListenAndServe(":"+cfg.Server.Port, r)
}
