package main

import (
	"fmt"
	"log"

	cacheserver "github.com/saikrishnamohan7/distributed-cache/cache_server"
	distributedcache "github.com/saikrishnamohan7/distributed-cache/distributed_cache"
)

func main() {
	cache := distributedcache.New()
	serverOptions := cacheserver.ServerOptions{
		ListenAddr: ":3000",
		IsLeader:   true,
	}

	server := cacheserver.New(serverOptions, cache)

	fmt.Println("Server running on", serverOptions.ListenAddr)
	err := server.Start()
	if err != nil {
		log.Fatal(err)
	}
}
