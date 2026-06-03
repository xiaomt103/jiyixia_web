package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"jiyixia/member-system/internal/auth"
	"jiyixia/member-system/internal/config"
	"jiyixia/member-system/internal/httpapi"
	"jiyixia/member-system/internal/store"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := store.NewPostgres(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect postgres: %v", err)
	}
	defer db.Close()
	if err := db.Migrate(ctx); err != nil {
		log.Fatalf("migrate postgres: %v", err)
	}

	sessions := auth.NewRedisSessionStore(cfg.RedisURL, cfg.SessionTTL)
	if err := sessions.Ping(ctx); err != nil {
		log.Fatalf("connect redis: %v", err)
	}
	defer sessions.Close()

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpapi.NewRouter(db, sessions, cfg),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("member api listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown: %v", err)
	}
}
