# Axinova MCP Server

**Model Context Protocol (MCP) v2025-11-25 Implementation**

Provides 38 DevOps tools across 5 services for MCP-native clients.

## 🚀 Quick Start

### Claude Desktop / Claude Code

```bash
# Install server
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-macos -o /usr/local/bin/axinova-mcp-server
chmod +x /usr/local/bin/axinova-mcp-server

# Configure Claude Desktop
# Edit: ~/Library/Application Support/Claude/claude_desktop_config.json
```

```json
{
  "mcpServers": {
    "axinova-tools": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
        "APP_PORTAINER__TOKEN": "ptr_YOUR_TOKEN_HERE",
        "APP_GRAFANA__URL": "https://grafana.axinova-internal.xyz",
        "APP_GRAFANA__TOKEN": "glsa_YOUR_TOKEN_HERE",
        "APP_TLS__SKIP_VERIFY": "true"
      }
    }
  }
}
```

### GitHub Copilot (VS Code)

```json
// settings.json
{
  "github.copilot.mcp.servers": {
    "axinova-tools": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
        "APP_PORTAINER__TOKEN": "ptr_YOUR_TOKEN_HERE",
        "APP_TLS__SKIP_VERIFY": "true"
      }
    }
  }
}
```

## 📚 Documentation

### Native MCP Integration (Recommended)

"Out of the box" integration with MCP-native clients:

- **[Native MCP Integration Guide](docs/NATIVE-MCP-INTEGRATION.md)** - Overview and architecture
- **[Claude Desktop Onboarding](docs/onboarding/claude-desktop.md)** - Step-by-step setup
- **[Claude Code Onboarding](docs/onboarding/claude-code.md)** - CLI integration
- **[GitHub Copilot Onboarding](docs/onboarding/github-copilot.md)** - VS Code/JetBrains setup

### API Integration (For Non-MCP Clients)

For platforms without native MCP support (ChatGPT, Gemini, LangChain, custom apps):

- **[API Reference](./docs/API-REFERENCE.md)** - Complete HTTP API documentation
- **[LLM Integration Guide](./docs/LLM-INTEGRATION.md)** - Integration examples
- **[Tool Catalog](./docs/TOOL-CATALOG.md)** - Complete list of all tools

**Production Server:** `https://mcp.axinova-ai.com`

## Overview

This MCP server exposes a unified interface for AI assistants and agents to interact with Axinova's internal infrastructure and productivity tools:

- **Portainer** (8 tools) - Docker container management
- **Grafana** (9 tools) - Monitoring dashboards and visualization
- **Prometheus** (7 tools) - Metrics queries and alerting
- **SilverBullet** (6 tools) - Note-taking and knowledge management
- **Vikunja** (8 tools) - Task and project management

## 🔌 Supported Clients

- ✅ **Claude Desktop** (native stdio)
- ✅ **Claude Code CLI** (native stdio)
- ✅ **GitHub Copilot** (VS Code, JetBrains, CLI)
- ✅ **Any MCP-compliant client**
- ✅ **ChatGPT, Gemini, LangChain** (via HTTP API)

## Features

- ✅ **MCP Protocol 2025-11-25** compliant
- ✅ **38 Tools** across 5 services
- ✅ **stdio transport** for local/container integrations
- ✅ **SSE transport** for remote web clients
- ✅ **HTTP JSON-RPC API** for non-MCP clients
- ✅ **Modular design** - easy to add new services
- ✅ **Configuration-driven** - Koanf with env var support
- ✅ **Production-ready** - Docker support, graceful shutdown, Prometheus metrics

## Documentation

📚 **Comprehensive documentation available in the [docs/](./docs) directory:**

- **[API Reference](./docs/API-REFERENCE.md)** - Complete HTTP API documentation with examples
- **[LLM Integration Guide](./docs/LLM-INTEGRATION.md)** - How to integrate with Claude, LangChain, LlamaIndex, OpenAI, and custom agents
- **[Tool Catalog](./docs/TOOL-CATALOG.md)** - Complete list of all 40+ tools with schemas and examples

**Production Server:** `https://mcp.axinova-ai.com`

## Quick Start

### Prerequisites

- Go 1.22+
- Access to Axinova internal tools (ax-tools, ax-sas-tools)
- API tokens for each service

### Installation

```bash
# Clone repository
git clone https://github.com/axinova-ai/axinova-mcp-server-go.git
cd axinova-mcp-server-go

# Install dependencies
go mod tidy

# Build
make build

# Or install globally
make install
```

### Configuration

Create a `.env` file from the example:

```bash
cp .env.example .env
# Edit .env with your tokens
```

Or use environment variables:

```bash
export APP_PORTAINER__URL=https://portainer.axinova-internal.xyz
export APP_PORTAINER__TOKEN=your-token
# ... etc
```

### Run

```bash
# Development
make run

# Production
ENV=prod ./bin/axinova-mcp-server
```

### Docker

```bash
# Build image
make docker-build

# Run with docker-compose
docker-compose up -d
```


## 🛠 Available Tools (38)

See [Tool Catalog](./docs/TOOL-CATALOG.md) for complete schemas and examples.

### Quick Reference

- **Portainer** (8) - `portainer_list_containers`, `portainer_start_container`, `portainer_stop_container`, `portainer_restart_container`, `portainer_get_container_logs`, `portainer_list_stacks`, `portainer_get_stack`, `portainer_get_container`
- **Grafana** (9) - `grafana_list_dashboards`, `grafana_get_dashboard`, `grafana_search_dashboards`, `grafana_get_dashboard_panels`, `grafana_list_datasources`, `grafana_list_folders`, `grafana_list_alerts`, `grafana_get_alert`, `grafana_test_datasource`
- **Prometheus** (7) - `prometheus_query`, `prometheus_query_range`, `prometheus_get_targets`, `prometheus_get_alerts`, `prometheus_get_alert_rules`, `prometheus_get_metrics`, `prometheus_get_label_values`
- **SilverBullet** (6) - `silverbullet_list_pages`, `silverbullet_read_page`, `silverbullet_search`, `silverbullet_get_page_meta`, `silverbullet_list_templates`, `silverbullet_query`
- **Vikunja** (8) - `vikunja_list_projects`, `vikunja_get_project`, `vikunja_list_tasks`, `vikunja_get_task`, `vikunja_create_task`, `vikunja_update_task`, `vikunja_list_labels`, `vikunja_search_tasks`

## 🌐 Integration Methods

### Native MCP Protocol (Recommended)

**Transport:** stdio or SSE
**Clients:** Claude Desktop, Claude Code, GitHub Copilot
**Setup:** [Native Integration Guide](docs/NATIVE-MCP-INTEGRATION.md)

Use this for the best experience with MCP-native clients. The server runs as a subprocess and communicates via standard input/output.

### HTTP API (For Non-MCP Clients)

**Transport:** HTTP JSON-RPC
**Clients:** ChatGPT, Gemini, LangChain, custom apps
**Setup:** [API Integration Guide](docs/UNIVERSAL-API-INTEGRATION.md)

Use this for platforms that don't support native MCP protocol.

## Usage Examples

### Natural Language Interaction

Once configured, interact using natural language:

**Container Management:**
```
You: Show me all running Docker containers
You: Get the last 100 lines of logs from the grafana container
You: Restart the postgres container
```

**Monitoring:**
```
You: List all Grafana dashboards
You: What's the current CPU usage across all services?
You: Show me active Prometheus alerts
```

**Task Management:**
```
You: Create a task "Review production logs" in the DevOps project
You: List all high-priority tasks
```

### Programmatic Usage (Embedding in Go Apps)

```go
import (
    "context"
    "github.com/axinova-ai/axinova-mcp-server-go/internal/mcp"
    "github.com/axinova-ai/axinova-mcp-server-go/internal/clients/portainer"
)

// Create MCP server
server := mcp.NewServer("my-app", "1.0.0", "2025-11-25")

// Register tools
portainerClient := portainer.NewClient(url, token, timeout)
portainer.RegisterTools(server, portainerClient)

// Run server
server.Run(context.Background())
```

## Architecture

```
cmd/
  server/              # Main entrypoint
internal/
  mcp/                 # MCP protocol implementation
    server.go          # JSON-RPC 2.0 stdio server
    types.go           # MCP type definitions
  clients/
    portainer/         # Portainer API client + tools
    grafana/           # Grafana API client + tools
    prometheus/        # Prometheus API client + tools
    silverbullet/      # SilverBullet API client + tools
    vikunja/           # Vikunja API client + tools
  config/              # Configuration management (Koanf)
config/
  base.yaml            # Base configuration
  dev.yaml             # Dev overrides
```

## Configuration

Configuration follows the Koanf convention with these precedence levels:

1. `config/base.yaml` (lowest)
2. `config/{ENV}.yaml`
3. Environment variables with `APP_` prefix (highest)

Example environment variables:

```bash
APP_PORTAINER__URL=https://portainer.example.com
APP_PORTAINER__TOKEN=secret
APP_LOG__LEVEL=debug
```

Note: Use double underscores (`__`) for nested keys.

## Development

```bash
# Run tests
make test

# Format code
make fmt

# Tidy dependencies
make tidy

# Clean build artifacts
make clean
```

## Adding New Services

1. Create client in `internal/clients/{service}/client.go`
2. Define tools in `internal/clients/{service}/tools.go`
3. Register in `cmd/server/main.go`
4. Add configuration to `config/base.yaml`
5. Update README

Example:

```go
// internal/clients/myservice/tools.go
func RegisterTools(server *mcp.Server, client *Client) {
    server.RegisterTool(mcp.Tool{
        Name: "myservice_list_items",
        Description: "List all items",
        InputSchema: mcp.InputSchema{Type: "object"},
    }, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
        return client.ListItems(ctx)
    })
}
```

## License

MIT

## Contributing

Contributions welcome! Please open an issue or PR.

## Support

For issues, please file a GitHub issue at https://github.com/axinova-ai/axinova-mcp-server-go/issues

---

**Sources:**
- [Model Context Protocol Specification](https://modelcontextprotocol.io/specification/2025-11-25)
- [MCP GitHub](https://github.com/modelcontextprotocol/modelcontextprotocol)

<!-- Testing CI/CD pipeline -->
