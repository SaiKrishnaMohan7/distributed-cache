package main

import (
	"log"
	"time"

	cacheserver "github.com/saikrishnamohan7/distributed-cache/cache_server"
	distributedcache "github.com/saikrishnamohan7/distributed-cache/distributed_cache"
)

func main() {
	cache := distributedcache.NewCache(time.Second) // TODO: read from env var
	serverOptions := cacheserver.ServerOptions{
		ListenAddr: ":3000", // TODO: pick this from the env OR default to this. Good use case for direnv
		IsLeader:   true,
	}

	server := cacheserver.NewServer(serverOptions, cache)

	log.Printf("Server running on: %s", serverOptions.ListenAddr)
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
