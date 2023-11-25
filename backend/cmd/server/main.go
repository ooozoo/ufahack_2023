package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"ufahack_2023/internal/config"
	"ufahack_2023/internal/delivery/http/handlers/auth/login"
	"ufahack_2023/internal/delivery/http/handlers/auth/register"
	mwAuth "ufahack_2023/internal/delivery/http/middleware/auth"
	mwLogger "ufahack_2023/internal/delivery/http/middleware/logger"
	"ufahack_2023/internal/service/auth"
	"ufahack_2023/internal/storage/postgres"
	"ufahack_2023/pkg/logger/sl"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("starting ufahack",
		slog.String("env", cfg.Env),
		slog.String("version", "v0.0.1"),
	)

	log.Debug("debug message are enabled")

	storage, err := postgres.New(cfg.Database)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	authService := auth.New(log, storage, storage, cfg.JWT.Secret, cfg.JWT.TTL)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.StripSlashes)
	router.Use(mwAuth.NewAuth(log, cfg.JWT.Secret, authService))

	router.Route("/api/auth", func(r chi.Router) {
		r.Post("/login", login.New(log, authService))
		r.Post("/register", register.New(log, authService))
	})

	log.Info(
		"starting server",
		slog.String("address", cfg.Server.Address),
		slog.Int("port", cfg.Server.Port),
	)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Address, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("failed to start server")
			}
		}
	}()

	log.Info("server started")

	<-done
	log.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", sl.Err(err))
		return
	}

	if err := storage.Close(); err != nil {
		log.Error("failed to close storage", sl.Err(err))
		return
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
