package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/axinova-ai/axinova-mcp-server-go/internal/clients/grafana"
	"github.com/axinova-ai/axinova-mcp-server-go/internal/clients/portainer"
	"github.com/axinova-ai/axinova-mcp-server-go/internal/clients/prometheus"
	"github.com/axinova-ai/axinova-mcp-server-go/internal/clients/silverbullet"
	"github.com/axinova-ai/axinova-mcp-server-go/internal/clients/vikunja"
	"github.com/axinova-ai/axinova-mcp-server-go/internal/config"
	"github.com/axinova-ai/axinova-mcp-server-go/internal/mcp"
)

func main() {
	// Get environment (default: dev)
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	// Load configuration
	cfg, err := config.Load(env)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Create MCP server
	mcpServer := mcp.NewServer(
		cfg.Server.Name,
		cfg.Server.Version,
		cfg.Server.ProtocolVersion,
	)

	// Register Portainer tools
	if cfg.Portainer.Enabled && cfg.Portainer.URL != "" {
		portainerClient := portainer.NewClient(
			cfg.Portainer.URL,
			cfg.Portainer.Token,
			cfg.Timeout.HTTP,
			cfg.TLS.SkipVerify,
		)
		portainer.RegisterTools(mcpServer, portainerClient)
		log.Printf("✓ Portainer tools registered (%s)", cfg.Portainer.URL)
	} else {
		log.Println("⊗ Portainer disabled or not configured")
	}

	// Register Grafana tools
	if cfg.Grafana.Enabled && cfg.Grafana.URL != "" {
		grafanaClient := grafana.NewClient(
			cfg.Grafana.URL,
			cfg.Grafana.Token,
			cfg.Timeout.HTTP,
			cfg.TLS.SkipVerify,
		)
		grafana.RegisterTools(mcpServer, grafanaClient)
		log.Printf("✓ Grafana tools registered (%s)", cfg.Grafana.URL)
	} else {
		log.Println("⊗ Grafana disabled or not configured")
	}

	// Register Prometheus tools
	if cfg.Prometheus.Enabled && cfg.Prometheus.URL != "" {
		prometheusClient := prometheus.NewClient(
			cfg.Prometheus.URL,
			cfg.Timeout.HTTP,
			cfg.TLS.SkipVerify,
		)
		prometheus.RegisterTools(mcpServer, prometheusClient)
		log.Printf("✓ Prometheus tools registered (%s)", cfg.Prometheus.URL)
	} else {
		log.Println("⊗ Prometheus disabled or not configured")
	}

	// Register SilverBullet tools
	if cfg.SilverBullet.Enabled && cfg.SilverBullet.URL != "" {
		silverbulletClient := silverbullet.NewClient(
			cfg.SilverBullet.URL,
			cfg.SilverBullet.Token,
			cfg.Timeout.HTTP,
			cfg.TLS.SkipVerify,
		)
		silverbullet.RegisterTools(mcpServer, silverbulletClient)
		log.Printf("✓ SilverBullet tools registered (%s)", cfg.SilverBullet.URL)
	} else {
		log.Println("⊗ SilverBullet disabled or not configured")
	}

	// Register Vikunja tools
	if cfg.Vikunja.Enabled && cfg.Vikunja.URL != "" {
		vikunjaClient := vikunja.NewClient(
			cfg.Vikunja.URL,
			cfg.Vikunja.Token,
			cfg.Timeout.HTTP,
			cfg.TLS.SkipVerify,
		)
		vikunja.RegisterTools(mcpServer, vikunjaClient)
		log.Printf("✓ Vikunja tools registered (%s)", cfg.Vikunja.URL)
	} else {
		log.Println("⊗ Vikunja disabled or not configured")
	}

	log.Println("========================================")
	log.Printf("MCP Server: %s v%s", cfg.Server.Name, cfg.Server.Version)
	log.Printf("Protocol: %s", cfg.Server.ProtocolVersion)
	log.Println("========================================")

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal, stopping...")
		cancel()
	}()

	// Run MCP server (stdio transport)
	log.Println("MCP Server starting (stdio transport)...")
	if err := mcpServer.Run(ctx); err != nil {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("MCP Server stopped")
}
