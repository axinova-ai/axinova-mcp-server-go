package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/axinova-ai/axinova-mcp-server-go/internal/mcp"
	"github.com/axinova-ai/axinova-mcp-server-go/internal/metrics"
)

// APIServer provides HTTP JSON-RPC interface to MCP server
type APIServer struct {
	port      int
	apiToken  string
	mcpServer *mcp.Server
	server    *http.Server
	logger    *log.Logger
}

// NewAPIServer creates a new API server
func NewAPIServer(port int, apiToken string, mcpServer *mcp.Server, logger *log.Logger) *APIServer {
	return &APIServer{
		port:      port,
		apiToken:  apiToken,
		mcpServer: mcpServer,
		logger:    logger,
	}
}

// Start starts the HTTP API server
func (a *APIServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// MCP JSON-RPC endpoint
	mux.HandleFunc("/api/mcp/v1/call", a.authMiddleware(a.handleRPCCall))

	// Tools list endpoint (convenience)
	mux.HandleFunc("/api/mcp/v1/tools", a.authMiddleware(a.handleListTools))

	a.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", a.port),
		Handler: mux,
	}

	// Shutdown on context cancellation
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		a.server.Shutdown(shutdownCtx)
	}()

	a.logger.Printf("MCP API server starting on port %d", a.port)
	return a.server.ListenAndServe()
}

// authMiddleware validates Bearer token
func (a *APIServer) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			a.sendError(w, http.StatusUnauthorized, -32000, "Missing Authorization header", nil)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.sendError(w, http.StatusUnauthorized, -32000, "Invalid Authorization header format", nil)
			return
		}

		token := parts[1]
		if token != a.apiToken {
			a.sendError(w, http.StatusUnauthorized, -32000, "Invalid API token", nil)
			return
		}

		metrics.ActiveConnections.Inc()
		defer metrics.ActiveConnections.Dec()

		next(w, r)
	}
}

// handleRPCCall handles JSON-RPC method calls
func (a *APIServer) handleRPCCall(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		a.sendError(w, http.StatusMethodNotAllowed, -32000, "Method not allowed", nil)
		return
	}

	startTime := time.Now()

	var req mcp.JSONRPCRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		a.sendError(w, http.StatusBadRequest, -32700, "Parse error", err.Error())
		return
	}

	// Handle the RPC request via MCP server
	result, err := a.mcpServer.HandleHTTPRequest(r.Context(), &req)
	duration := time.Since(startTime)

	if err != nil {
		metrics.RecordRPCRequest(req.Method, "http", duration, "error")
		a.sendError(w, http.StatusInternalServerError, -32603, "Internal error", err.Error())
		return
	}

	metrics.RecordRPCRequest(req.Method, "http", duration, "")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      req.ID,
		"result":  result,
	})
}

// handleListTools returns available tools
func (a *APIServer) handleListTools(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		a.sendError(w, http.StatusMethodNotAllowed, -32000, "Method not allowed", nil)
		return
	}

	tools := a.mcpServer.GetTools()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tools": tools,
		"count": len(tools),
	})
}

// sendError sends JSON-RPC error response
func (a *APIServer) sendError(w http.ResponseWriter, httpStatus, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"jsonrpc": "2.0",
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"data":    data,
		},
		"id": nil,
	})
}
