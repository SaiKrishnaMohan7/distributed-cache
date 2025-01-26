package cacheserver

import (
	"io"
	"net/http"
	"time"

	distributedcache "github.com/saikrishnamohan7/distributed-cache/distributed_cache"
)

type ServerOptions struct {
	ListenAddr string
	IsLeader   bool
}

type Server struct {
	options ServerOptions
	cache   distributedcache.Cache
}

func New(options ServerOptions, cache *distributedcache.InMemoryCache) *Server {
	return &Server{
		options: options,
		cache:   cache,
	}
}

func (server *Server) Start() error {
	http.HandleFunc("/get", server.handleGetCommand)       // HTTP GET
	http.HandleFunc("/set", server.handleSetCommand)       // HTTP POST
	http.HandleFunc("/has", server.handleHasCommand)       // HTTP GET
	http.HandleFunc("/delete", server.handleDeleteCommand) // HTTP DELETE

	return http.ListenAndServe(":3000", nil)
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

	ttl := 0 * time.Second
	if ttlStr == "" {
		ttl = 0
	} else {
		parsedTTL, err := time.ParseDuration(ttlStr)
		if err != nil {
			http.Error(responseWriter, "Invalid TTL", http.StatusBadRequest)
			return
		}

		ttl = parsedTTL
	}

	buf, err := io.ReadAll(value)
	if err != nil {
		http.Error(responseWriter, "Failed to read body", http.StatusInternalServerError)
		return
	}

	server.cache.Set([]byte(key), buf, ttl)
	responseWriter.WriteHeader(http.StatusOK)
}

func (server *Server) handleHasCommand(responseWriter http.ResponseWriter, req *http.Request) {
	key := req.URL.Query().Get("key")

	if key == "" {
		http.Error(responseWriter, "Invalid Key", http.StatusBadRequest)
		return
	}

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

	keyByte := []byte(key)

	server.cache.Delete(keyByte)

	responseWriter.WriteHeader(http.StatusOK)
}
