package cacheserver

import (
	"context"
	"io"
	"log"
	"net/http"
	"time"

	distcache "github.com/saikrishnamohan7/distributed-cache/internal/cache"
)

// Server is
type Server struct {
	cache      *distcache.InMemoryCache
	httpServer *http.Server
}

// NewServer is a constructor fn
func NewServer(cache *distcache.InMemoryCache, port string) *Server {
	mux := http.NewServeMux()
	server := &Server{
		cache: cache,
	}

	// Routing
	mux.HandleFunc("/get", server.handleGetCommand)       // HTTP GET
	mux.HandleFunc("/set", server.handleSetCommand)       // HTTP POST
	mux.HandleFunc("/has", server.handleHasCommand)       // HTTP GET
	mux.HandleFunc("/delete", server.handleDeleteCommand) // HTTP DELETE

	server.httpServer = &http.Server{
		Addr:    port,
		Handler: mux,
	}

	return server
}

// Start starts the server
func (server *Server) Start() error {
	log.Println("Starting server on ", server.httpServer.Addr)
	return server.httpServer.ListenAndServe()
}

// Shutdown handles graceful wind down of server
func (server *Server) Shutdown(ctx context.Context) error {
	log.Println("Gracefully shutting down server")
	log.Println("Signalling Cache Clean up stoppage")
	server.cache.StopCleanup()

	return server.httpServer.Shutdown(ctx)
}

func (server *Server) handleGetCommand(responseWriter http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	if key == "" {
		http.Error(responseWriter, "Invalid Key", http.StatusBadRequest)
		return
	}

	value, err := server.cache.Get([]byte(key))

	if err != nil {
		http.Error(responseWriter, "Not Found", http.StatusNotFound)
		return
	}

	log.Printf("Received GET request for key: %s", key)

	responseWriter.Write(value)
}

func (server *Server) handleSetCommand(responseWriter http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")
	ttlStr := req.URL.Query().Get("ttl")

	if key == "" {
		http.Error(responseWriter, "Invalid Key", http.StatusBadRequest)
		return
	}

	value := req.Body
	defer value.Close()

	parsedTTL, err := time.ParseDuration(ttlStr)
	if err != nil {
		http.Error(responseWriter, "Invalid TTL", http.StatusBadRequest)
		return
	}

	buf, err := io.ReadAll(value)
	if err != nil {
		http.Error(responseWriter, "Failed to read body", http.StatusInternalServerError)
		return
	}

	log.Printf("Received SET request with Key: %s, value: %s and TTL: %s", key, buf, ttlStr)

	if err := server.cache.Set([]byte(key), buf, parsedTTL); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	responseWriter.WriteHeader(http.StatusOK)
}

func (server *Server) handleHasCommand(responseWriter http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	if key == "" {
		http.Error(responseWriter, "Invalid Key", http.StatusBadRequest)
		return
	}

	log.Printf("Received, HAS request for key: %s", key)
	hasKey := server.cache.Has([]byte(key))

	if hasKey {
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte("true"))
	} else {
		responseWriter.WriteHeader(http.StatusNotFound)
		responseWriter.Write([]byte("false"))
	}

}

func (server *Server) handleDeleteCommand(responseWriter http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	if key == "" {
		http.Error(responseWriter, "Invalid Key", http.StatusBadRequest)
		return
	}
	log.Printf("Received, DELETE request for key: %s", key)

	keyByte := []byte(key)

	server.cache.Delete(keyByte)

	responseWriter.WriteHeader(http.StatusOK)
}
