package main

import (
	"LostAndFound/internal/adapters/postgres"
	myredis "LostAndFound/internal/adapters/redis"
	mys3 "LostAndFound/internal/adapters/s3"
	"LostAndFound/internal/auth"
	"LostAndFound/internal/bootstrap"
	server_config "LostAndFound/internal/config/server_config"
	storage_config "LostAndFound/internal/config/storage_config"
	router "LostAndFound/internal/delivery/http"
	"LostAndFound/internal/delivery/http/handler"
	"LostAndFound/internal/service"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file")
		os.Exit(1)
	}
}

// @title           LostAndFound API
// @version         1.0
// @description     API для поиска и возврата потерянных вещей
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {

	serverCfg, err := server_config.MustLoadServerConfig()
	if err != nil {
		slog.Error("Error loading server config")
		os.Exit(1)
	}

	storageCfg, err := storage_config.MustLoadStorageConfig()
	if err != nil {
		slog.Error("Error loading storage config")
		os.Exit(1)
	}

	postgresDb, err := postgres.NewStorage(storageCfg.Postgres)
	if err != nil {
		slog.Error("failed to connect to postgres", err)
		os.Exit(1)
	}

	redis, err := myredis.NewRedis(storageCfg.Redis)
	if err != nil {
		slog.Error("failed to connect to redis", err)
		os.Exit(1)
	}

	s3, err := mys3.NewS3Client(storageCfg.S3)
	if err != nil {
		slog.Error("failed to connect to s3", err)
		os.Exit(1)
	}

	slog.Info("connected to database")

	defer func() {
		err = postgresDb.Close()
		if err != nil {
			slog.Error("got error when closing the DB connection", err)
			os.Exit(1)
		}
	}()

	repos := bootstrap.Init(postgresDb, redis, s3, storageCfg.S3)

	tokenManager, err := auth.NewTokenManager(repos.CacheRepo)
	if err != nil {
		slog.Error("failed to initialize token manager", err)
		os.Exit(1)
	}

	services := service.NewService(repos, tokenManager)

	handlers := handler.NewHandler(services, tokenManager)

	slog.Info("starting server")

	server := &http.Server{
		Addr:    serverCfg.Address,
		Handler: router.NewRouter(handlers, redis),
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("starting server...")
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("listen and serve error", err)
		}
	}()

	<-quit
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", err)
	}

	slog.Info("server exiting")
}
