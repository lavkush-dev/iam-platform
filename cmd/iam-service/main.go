package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"database/sql"

	"iam-platform/internal/config"
	"iam-platform/internal/handlers"
	"iam-platform/internal/logger"
	appmw "iam-platform/internal/middleware"
	"iam-platform/internal/repository/postgres"
	"iam-platform/internal/services"
)

func main() {
	cfg := config.Load()

	log := newLogger()
	defer log.Sync()

	db := postgres.NewDB(cfg.DB.URL)

	router := setupRouter(db, log)

	server := newHTTPServer(cfg.Server.Port, router)

	runServer(server, log)

	gracefulShutdown(server, log)
}

func newLogger() *zap.Logger {
	return logger.New()
}

func setupRouter(db *sql.DB, log *zap.Logger) http.Handler {
	userRepo := postgres.NewUserRepository(db)

	userService := services.NewUserService(userRepo)

	userHandler := handlers.NewUserHandler(userService, log)

	r := chi.NewRouter()

	registerMiddleware(r, log)

	registerHealthRoutes(r)
	registerUserRoutes(r, userHandler)

	return r
}

func registerMiddleware(
	r chi.Router,
	log *zap.Logger,
) {
	r.Use(appmw.ZapLogger(log))
}

func registerHealthRoutes(r chi.Router) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

func registerUserRoutes(
	r chi.Router,
	userHandler *handlers.UserHandler,
) {
	r.Post("/users", userHandler.CreateUser)

	r.Get("/users/{id}", userHandler.GetUser)
}

func newHTTPServer(
	port string,
	handler http.Handler,
) *http.Server {
	return &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
}

func runServer(
	server *http.Server,
	logger *zap.Logger,
) {
	go func() {
		logger.Info(
			"iam-service running",
			zap.String("addr", server.Addr),
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			logger.Fatal(
				"failed to start server",
				zap.Error(err),
			)
		}
	}()
}

func gracefulShutdown(
	server *http.Server,
	log *zap.Logger,
) {
	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-quit

	log.Info("shutdown signal received")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {

		log.Error(
			"graceful shutdown failed",
			zap.Error(err),
		)

		if err := server.Close(); err != nil {
			log.Error(
				"force close failed",
				zap.Error(err),
			)
		}
	}

	log.Info("server exited cleanly")
}
