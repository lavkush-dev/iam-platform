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
	"iam-platform/pkg/jwt"
)

func main() {
	cfg := config.Load()

	log := newLogger()
	defer log.Sync()

	db := postgres.NewDB(cfg.DB.URL)

	router := setupRouter(db, cfg, log)

	server := newHTTPServer(
		cfg.Server.Port,
		router,
	)

	startServer(server, log)

	gracefulShutdown(server, log)
}

func newLogger() *zap.Logger {
	return logger.New()
}

func setupRouter(
	db *sql.DB,
	cfg *config.Config,
	log *zap.Logger,
) http.Handler {

	// repositories
	userRepo := postgres.NewUserRepository(db)

	// external dependencies
	jwtManager := jwt.NewManager(cfg.JWT.Secret)

	// services
	authService := services.NewAuthService(
		userRepo,
		jwtManager,
	)

	// handlers
	authHandler := handlers.NewAuthHandler(
		authService,
		log,
	)

	r := chi.NewRouter()

	registerMiddleware(r, log)

	registerHealthRoutes(r)

	registerAuthRoutes(r, authHandler)

	return r
}

func registerMiddleware(
	r chi.Router,
	log *zap.Logger,
) {
	r.Use(appmw.ZapLogger(log))
}

func registerHealthRoutes(r chi.Router) {
	r.Get("/health", func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		w.Header().Set(
			"Content-Type",
			"application/json",
		)

		w.WriteHeader(http.StatusOK)

		w.Write([]byte("OK"))
	})
}

func registerAuthRoutes(
	r chi.Router,
	authHandler *handlers.AuthHandler,
) {
	r.Post("/login", authHandler.Login)
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

func startServer(
	server *http.Server,
	log *zap.Logger,
) {
	go func() {

		log.Info(
			"auth-service running",
			zap.String("addr", server.Addr),
		)

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {

			log.Fatal(
				"server failed",
				zap.Error(err),
			)
		}
	}()
}

func gracefulShutdown(
	server *http.Server,
	log *zap.Logger,
) {
	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

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

		if closeErr := server.Close(); closeErr != nil {

			log.Error(
				"force close failed",
				zap.Error(closeErr),
			)
		}
	}

	log.Info("server exited cleanly")
}
