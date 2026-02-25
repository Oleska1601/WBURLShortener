package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Oleska1601/WBURLShortener/config"
	api "github.com/Oleska1601/WBURLShortener/internal/controller"
	v1 "github.com/Oleska1601/WBURLShortener/internal/controller/api/v1"
	"github.com/Oleska1601/WBURLShortener/internal/redis"
	"github.com/Oleska1601/WBURLShortener/internal/repo/postgres"
	"github.com/Oleska1601/WBURLShortener/internal/service"
	"github.com/wb-go/wbf/zlog"
)

// @title URL Shortener
// @version 1.0
// @description API for URL Shortener
// @termsOfService http://swagger.io/terms/

// @host localhost:8081
// @BasePath /
func Run(cfg *config.Config) {
	zlog.Init()
	if err := zlog.SetLevel(cfg.Logger.Level); err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run zlog.SetLevel").
			Msg("failed to set log level")
	}

	db, err := initDB(&cfg.Database.Postgres)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run initDB").
			Msg("failed to init database")
	}

	repo := postgres.New(db)
	if err := repo.ApplyMigrations(); err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run repo.ApplyMigrations").
			Msg("failed to apply migrations")
	}

	redis, err := redis.New(&cfg.Redis)
	if err != nil {
		zlog.Logger.Fatal().
			Err(err).
			Str("path", "Run redis.New").
			Msg("init redis")
	}

	service := service.New(redis, repo)
	apiV1 := v1.New(service)
	router := api.Register(&cfg.Gin, apiV1)
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{Addr: addr, Handler: router}

	go func() {
		zlog.Logger.Info().Str("path", "Run").Str("addr", addr).Msg("start server")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zlog.Logger.Fatal().
				Err(err).
				Str("path", "Run server.ListenAndServe").
				Msg("failed to process server")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		zlog.Logger.Err(err).Str("path", "App server.Shutdown").
			Msg("failed to shutdown server")
	}

	zlog.Logger.Info().Msg("shutdown server properly")
}
