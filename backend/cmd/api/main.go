package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kanvix/backend/internal/cache"
	"kanvix/backend/internal/config"
	"kanvix/backend/internal/database"
	"kanvix/backend/internal/http/routes"
	"kanvix/backend/internal/logging"
	"kanvix/backend/internal/repositories/gormrepo"
	"kanvix/backend/internal/services"
)

func main() {
	log := logging.New()
	slog.SetDefault(log)

	cfg := config.Load()
	ctx := context.Background()

	db, err := database.Open(ctx, log, cfg.DBPath)
	if err != nil {
		log.Error("db open failed", "error", err)
		os.Exit(1)
	}
	if err := database.AutoMigrate(ctx, db.Gorm); err != nil {
		log.Error("db migrate failed", "error", err)
		os.Exit(1)
	}
	if err := database.EnsureSeeded(ctx, db.Gorm); err != nil {
		log.Error("db seed failed", "error", err)
		os.Exit(1)
	}

	redisWrap := cache.NewRedis(cfg.RedisAddr, cfg.RedisPassword, cfg.RedisDB)
	if err := redisWrap.Ping(ctx); err != nil {
		log.Error("redis unavailable", "error", err)
		os.Exit(1)
	}

	repo := gormrepo.New(db.Gorm)
	authSvc := services.AuthService{Repo: repo, Redis: redisWrap.Client, SessionTTL: cfg.SessionTTL}
	appSvc := services.AppService{Repo: repo, Redis: redisWrap.Client, CacheTTL: cfg.CacheTTL}

	r := routes.NewRouter(routes.Deps{
		Log:        log,
		Cfg:        cfg,
		Auth:       authSvc,
		AppService: appSvc,
	})

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("server listening", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("server error", "error", err)
		}
	}()

	<-stop
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = redisWrap.Client.Close()
	_ = db.SQL.Close()
	_ = srv.Shutdown(shutdownCtx)
}
