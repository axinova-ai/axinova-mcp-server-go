# Axinova MCP Server

Model Context Protocol (MCP) server providing LLM/agent access to Axinova internal tools.

## Overview

This MCP server exposes a unified interface for AI assistants and agents to interact with Axinova's internal infrastructure and productivity tools:

- **Portainer** - Docker container management
- **Grafana** - Monitoring dashboards and visualization
- **Prometheus** - Metrics queries and alerting
- **SilverBullet** - Note-taking and knowledge management
- **Vikunja** - Task and project management

## Features

- ✅ **MCP Protocol 2025-11-25** compliant
- ✅ **40+ Tools** across 5 services
- ✅ **stdio transport** for local/container integrations
- ✅ **HTTP JSON-RPC API** for remote LLM agent access
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

## Available Tools

### Portainer Tools

| Tool | Description |
|------|-------------|
| `portainer_list_containers` | List all containers |
| `portainer_start_container` | Start a container |
| `portainer_stop_container` | Stop a container |
| `portainer_restart_container` | Restart a container |
| `portainer_get_container_logs` | Get container logs |
| `portainer_list_stacks` | List Docker Compose stacks |
| `portainer_get_stack` | Get stack details |
| `portainer_inspect_container` | Inspect container |

### Grafana Tools

| Tool | Description |
|------|-------------|
| `grafana_list_dashboards` | List all dashboards |
| `grafana_get_dashboard` | Get dashboard by UID |
| `grafana_create_dashboard` | Create new dashboard |
| `grafana_delete_dashboard` | Delete dashboard |
| `grafana_list_datasources` | List datasources |
| `grafana_create_datasource` | Create datasource |
| `grafana_query_datasource` | Query datasource (PromQL, etc.) |
| `grafana_list_alert_rules` | List alert rules |
| `grafana_get_health` | Check Grafana health |

### Prometheus Tools

| Tool | Description |
|------|-------------|
| `prometheus_query` | Execute instant query |
| `prometheus_query_range` | Execute range query |
| `prometheus_list_label_names` | Get all label names |
| `prometheus_list_label_values` | Get label values |
| `prometheus_find_series` | Find time series |
| `prometheus_list_targets` | List scrape targets |
| `prometheus_get_metadata` | Get metric metadata |

### SilverBullet Tools

| Tool | Description |
|------|-------------|
| `silverbullet_list_pages` | List all pages |
| `silverbullet_get_page` | Get page content |
| `silverbullet_create_page` | Create new page |
| `silverbullet_update_page` | Update page |
| `silverbullet_delete_page` | Delete page |
| `silverbullet_search_pages` | Search pages |

### Vikunja Tools

| Tool | Description |
|------|-------------|
| `vikunja_list_projects` | List all projects |
| `vikunja_get_project` | Get project details |
| `vikunja_create_project` | Create project |
| `vikunja_list_tasks` | List tasks in project |
| `vikunja_get_task` | Get task details |
| `vikunja_create_task` | Create new task |
| `vikunja_update_task` | Update task |
| `vikunja_delete_task` | Delete task |

## Usage Examples

### With Claude Desktop

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "axinova": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "APP_PORTAINER__TOKEN": "your-token",
        "APP_GRAFANA__TOKEN": "your-token",
        "APP_VIKUNJA__TOKEN": "your-token",
        "APP_SILVERBULLET__TOKEN": "your-token"
      }
    }
  }
}
```

### Programmatic Usage

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
