package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HealthServer struct {
	port   int
	server *http.Server
}

func NewHealthServer(port int) *HealthServer {
	return &HealthServer{port: port}
}

func (hs *HealthServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// Health check endpoint (required for Docker health check)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"service":   "axinova-mcp-server",
		})
	})

	// Readiness check (can be extended to check service connectivity)
	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
	})

	// Status endpoint (server info)
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"version":  "0.1.0",
			"protocol": "2025-11-25",
			"mode":     "stdio",
		})
	})

	hs.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", hs.port),
		Handler: mux,
	}

	// Shutdown on context cancellation
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		hs.server.Shutdown(shutdownCtx)
	}()

	return hs.server.ListenAndServe()
}
