package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/saikrishnamohan7/distributed-cache/config"
	districache "github.com/saikrishnamohan7/distributed-cache/internal/cache" // should match go.mod, the "github.com/saikrishnamohan7/distributed-cache.. part. If go.mod changes, this changes too
	cacheserver "github.com/saikrishnamohan7/distributed-cache/internal/server"
)

func main() {
	config.InitLogger()
	err := config.LoadDotEnv("./.env")
	if err != nil {
		log.Fatalf("Unexpected error loading env vars: %v", err)
	}

	tickMs := 5000
	tickMsStr := os.Getenv("CACHE_CLEANUP_TICK")
	if tickMsStr != "" {
		var err error
		tickMs, err = strconv.Atoi(tickMsStr)
		if err != nil {
			log.Fatalf("invalid CACHE_CLEANUP_TICK: %v", err)
		}
	}
	if err != nil {
		log.Fatalf("invalid CACHE_CLEANUP_TICK: %v", err)
	}

	cleanupTick := time.Duration(tickMs) * time.Millisecond
	cache := districache.NewCache(cleanupTick)

	port := os.Getenv("PORT")
	if port == "" {
		port = ":42069"
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	cache.StartCleanup()
	server := cacheserver.NewServer(cache, port)

	// Start server in go routine so that I can handle graceful shutdown
	go func() {
		err = server.Start()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop // Block! Wait for SIGINT or SIGTERM

	log.Println("Shutting Down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown failed: %v", err)
	}
}
