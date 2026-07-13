package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"codelife-study-be/internal/adapter/cache"
	"codelife-study-be/internal/adapter/database"
	"codelife-study-be/internal/adapter/email"
	authrepo "codelife-study-be/internal/adapter/repository/auth"
	documentrepo "codelife-study-be/internal/adapter/repository/document"
	"codelife-study-be/internal/config"
	"codelife-study-be/internal/delivery/httpapi"
	domaindocument "codelife-study-be/internal/domain/document"
	authusecase "codelife-study-be/internal/usecase/auth"
	documentusecase "codelife-study-be/internal/usecase/document"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	postgres, err := database.OpenPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		logger.Error("postgres startup failed", "error", err)
		os.Exit(1)
	}
	if postgres != nil {
		defer postgres.Close()
		if err := database.RunMigrations(ctx, postgres); err != nil {
			logger.Error("postgres migration failed", "error", err)
			os.Exit(1)
		}
	}
	redis := cache.NewRedis(cfg.RedisAddress, cfg.RedisPassword, cfg.RedisDB, cfg.CacheTTL)
	var documentCache domaindocument.Cache
	var redisPinger httpapi.Pinger
	if redis != nil {
		documentCache = redis
		redisPinger = redis
	}
	var postgresPinger httpapi.Pinger
	if postgres != nil {
		postgresPinger = postgres
	}

	repository := documentrepo.NewEmbeddedRepository()
	documents := documentusecase.New(repository, documentCache)
	var authService *authusecase.Service
	if postgres != nil {
		authRepository := authrepo.NewPostgresRepository(postgres)
		mailer := email.NewSMTPMailer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUsername, cfg.SMTPPassword, cfg.SMTPFrom, logger)
		authService = authusecase.New(authRepository, mailer, cfg.AuthTokenSecret, cfg.AuthOTPTTL, cfg.AuthTokenTTL)
	}
	handler := httpapi.New(documents, authService, postgresPinger, redisPinger, logger, cfg.MaxBodyBytes)
	server := &http.Server{Addr: cfg.Address, Handler: handler, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 10 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}

	go func() {
		logger.Info("api started", "address", cfg.Address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server failed", "error", err)
			stop()
		}
	}()
	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
	}
}
