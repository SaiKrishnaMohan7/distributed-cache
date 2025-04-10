package main

import (
	"log"
	"time"

	"github.com/saikrishnamohan7/distributed-cache/internal/p2p"
)
func main() {
	node, err := p2p.NewNode( 2*time.Second, ":4000")
	if err != nil {
		log.Fatalf("Failed to start node: %v", err)
	}

	go node.Start()

	time.Sleep(10 * time.Second)
	node.Stop()
}