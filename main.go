package main

import (
	"log"
	"os"
	"strconv"
	"time"

	cacheserver "github.com/saikrishnamohan7/distributed-cache/cache_server"
	"github.com/saikrishnamohan7/distributed-cache/config"
	distributedcache "github.com/saikrishnamohan7/distributed-cache/distributed_cache"
)

func main() {
	err := config.LoadDotEnv("./.env")
	if err != nil {
		log.Fatalf("Unexpected error loading env vars: %v", err)
	}

	tickMs, err :=  strconv.Atoi(os.Getenv("CACHE_CLEANUP_TICK"))
	if err != nil {
		log.Fatalf("invalid CACHE_CLEANUP_TICK: %v", err)
	}

	cleanupTick := time.Duration(tickMs) * time.Millisecond
	cache := distributedcache.NewCache(cleanupTick)
	serverOptions := cacheserver.ServerOptions{
		ListenAddr: os.Getenv("PORT"),
		IsLeader:   true,
	}

	server := cacheserver.NewServer(serverOptions, cache)
	log.Printf("Server listening on: %s", serverOptions.ListenAddr)

	err = server.Start()
	if err != nil {
		log.Fatal(err)
	}
}


