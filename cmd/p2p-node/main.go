package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saikrishnamohan7/distributed-cache/internal/p2p"
)
func main() {
	node, err := p2p.NewNode( 2*time.Second, ":4000")
	if err != nil {
		log.Fatalf("Failed to start node: %v", err)
	}

	go node.Start()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop // wait for SIGINT or SIGTERM

	log.Println("Propagating shutdown signal...")
	node.Stop()
}