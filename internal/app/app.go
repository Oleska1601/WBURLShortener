package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Oleska1601/WBURLShortener/config"
	"github.com/Oleska1601/WBURLShortener/internal/controller"
	"github.com/Oleska1601/WBURLShortener/internal/database/repo"
	"github.com/Oleska1601/WBURLShortener/internal/redis"
	"github.com/Oleska1601/WBURLShortener/internal/usecase"
	"github.com/wb-go/wbf/zlog"
)

// @title URL Shortener
// @version 1.0
// @description API for URL Shortener
// @termsOfService http://swagger.io/terms/

// @host localhost:8082
// @BasePath /
func Run(cfg *config.Config) {
	// logger
	zlog.Init()
	if err := zlog.SetLevel(cfg.Logger.Level); err != nil {
		log.Fatalln("set zlog level error: %w", err)
	}

	// postgres
	db, err := initDB(&cfg.DB)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run initDB").
			Msg("init database")
	}

	// postgres repo
	pgRepo := repo.New(db)
	if err := pgRepo.ApplyMigrations(); err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run pgRepo.ApplyMigrations").
			Msg("apply migrations to database")
	}
	redis, err := redis.New(&cfg.Redis)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run redis.New").
			Msg("init redis")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	usecase := usecase.New(redis, pgRepo)
	server := controller.New(&cfg.Server, usecase)

	go func() {
		if err := server.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zlog.Logger.
				Fatal().
				Err(err).
				Str("path", "Run server.Srv.ListenAndServe").
				Msg("cannot start server")
		}
		zlog.Logger.Info().Msgf("server is started http://%s:%d/", cfg.Server.Host, cfg.Server.Port)
	}()

	// Ожидание сигнала для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zlog.Logger.Info().Msg("shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, cfg.Server.ShutdownTimeout)
	defer shutdownCancel()
	if err := server.Srv.Shutdown(shutdownCtx); err != nil {
		zlog.Logger.Error().Err(err).Msg("server shutdown")
		return
	}

	zlog.Logger.Info().Msg("server exited properly")

}
