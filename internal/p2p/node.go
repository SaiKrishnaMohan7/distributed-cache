package p2p

import (
	"log"
	"net"
	"time"

	"github.com/google/uuid"

	districache "github.com/saikrishnamohan7/distributed-cache/internal/cache"
)

type Node struct {
	listener net.Listener // 8 bytes pointer
	cache   *districache.InMemoryCache // 8 bytes pointer
	stop chan struct{} // 8 bytes pointer
	peers   []*Node // like a graph node; 8 bytes slice of pointers
	id      string // 16 bytes
	address string // this node's own address: ip:port; 16 bytes
}

func NewNode(cleanupTick time.Duration, address string) (*Node, error){
	listener, err := net.Listen("tcp", address);
	if err != nil {
		log.Fatalf("Error setting up TCP listener: %v", err)
		return nil, err
	}

	node := &Node{
		listener: listener,
		cache: districache.NewCache(cleanupTick),
		stop: make(chan struct{}),
		id: uuid.New().String(),
		address: address,
	}
	return node, nil
}

func (node *Node) Start() {
	log.Printf("Node %s listening on %s \n", node.id, node.address)

	for {
		conn, err := node.listener.Accept()
		if err != nil {
			select {
			case <-node.stop:
				log.Println("Node Shutting down...")
				return;
			default:
				log.Printf("Accept Error: %v", err)
				continue
			}
		}

		log.Println("Accepted connection")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("Hello!"))
}

func (node *Node) Stop() {
	log.Printf("Stop received, Stopping...")
	close(node.stop)
	node.listener.Close()
}