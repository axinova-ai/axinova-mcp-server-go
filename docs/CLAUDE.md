# CLAUDE.md

This file provides guidance to Claude Code when working with the axinova-mcp-server-go repository.

## Project Overview

**axinova-mcp-server-go** is a Model Context Protocol (MCP) server implementation that provides unified LLM/agent access to Axinova's internal infrastructure and productivity tools.

### Purpose

- Expose Portainer, Grafana, Prometheus, SilverBullet, and Vikunja via MCP protocol
- Enable AI assistants to manage containers, query metrics, create notes, and manage tasks
- Provide standardized tool interface for internal automation

### Architecture

```
┌─────────────────────────────────────┐
│         MCP Server (Go)             │
│  ┌──────────────────────────────┐   │
│  │   JSON-RPC 2.0 over stdio    │   │
│  └──────────────────────────────┘   │
│  ┌────┬────────┬──────┬─────────┐   │
│  │ P  │   G    │  Pr  │   S/V   │   │
│  │ o  │   r    │  o   │         │   │
│  │ r  │   a    │  m   │         │   │
│  │ t  │   f    │  e   │         │   │
│  │ a  │   a    │  t   │         │   │
│  │ i  │   n    │  h   │         │   │
│  │ n  │   a    │  e   │         │   │
│  │ e  │        │  u   │         │   │
│  │ r  │        │  s   │         │   │
│  └─┬──┴───┬────┴──┬───┴────┬────┘   │
└────┼──────┼───────┼────────┼────────┘
     │      │       │        │
     ▼      ▼       ▼        ▼
  Docker  Metrics  Time   Notes/Tasks
```

## Build and Development Commands

### Quick Reference

```bash
# Build
make build              # Build binary to bin/axinova-mcp-server
make install            # Install to /usr/local/bin

# Development
make run                # Run server directly
make test               # Run tests
make fmt                # Format code
make tidy               # Tidy dependencies

# Docker
make docker-build       # Build Docker image
make docker-push        # Push to ghcr.io/axinova-ai

# Testing
./test_mcp.sh           # Test MCP protocol
./test_tool_call.sh     # Test tool execution
```

### Configuration

```bash
# Environment selection
ENV=dev|prod           # Loads config/dev.yaml or config/prod.yaml

# Override via environment variables (APP_ prefix, __ for nesting)
export APP_PORTAINER__URL=https://portainer.axinova-internal.xyz
export APP_PORTAINER__TOKEN=your-token
export APP_TLS__SKIP_VERIFY=true
```

## Repository Structure

```
axinova-mcp-server-go/
├── cmd/
│   └── server/              # Main entrypoint (MCP server)
├── internal/
│   ├── mcp/
│   │   ├── server.go        # JSON-RPC 2.0 stdio server
│   │   └── types.go         # MCP type definitions
│   ├── clients/             # Service-specific API clients
│   │   ├── portainer/
│   │   │   ├── client.go    # Portainer API wrapper
│   │   │   └── tools.go     # MCP tool registrations
│   │   ├── grafana/         # Grafana client
│   │   ├── prometheus/      # Prometheus client
│   │   ├── silverbullet/    # SilverBullet client
│   │   └── vikunja/         # Vikunja client
│   └── config/
│       └── config.go        # Koanf configuration loader
├── config/
│   ├── base.yaml            # Base config (defaults)
│   ├── dev.yaml             # Development overrides
│   └── prod.yaml            # Production overrides
├── scripts/
│   └── get-tokens.md        # Token generation guide
├── Makefile                 # Build automation
├── Dockerfile               # Multi-stage Docker build
├── docker-compose.yml       # Deployment configuration
├── README.md                # User documentation
├── TESTING.md               # Testing guide
└── CLAUDE.md                # This file
```

## Technology Stack

- **Language**: Go 1.22+
- **Protocol**: MCP 2025-11-25 (JSON-RPC 2.0)
- **Configuration**: Koanf (YAML + env vars)
- **Transport**: stdio (standard input/output)
- **Deployment**: Docker, docker-compose

## Common Development Tasks

### Adding a New Service Client

1. **Create client package**
   ```bash
   mkdir -p internal/clients/newservice
   ```

2. **Implement client.go**
   ```go
   package newservice

   import (
       "context"
       "crypto/tls"
       "net/http"
       "time"
   )

   type Client struct {
       baseURL    string
       token      string
       httpClient *http.Client
   }

   func NewClient(baseURL, token string, timeout time.Duration, skipTLSVerify bool) *Client {
       transport := &http.Transport{
           TLSClientConfig: &tls.Config{
               InsecureSkipVerify: skipTLSVerify,
           },
       }

       return &Client{
           baseURL: baseURL,
           token:   token,
           httpClient: &http.Client{
               Timeout:   timeout,
               Transport: transport,
           },
       }
   }

   // Implement API methods
   func (c *Client) DoSomething(ctx context.Context) (interface{}, error) {
       // Implementation
   }
   ```

3. **Implement tools.go**
   ```go
   package newservice

   import (
       "context"
       "github.com/axinova-ai/axinova-mcp-server-go/internal/mcp"
   )

   func RegisterTools(server *mcp.Server, client *Client) {
       server.RegisterTool(mcp.Tool{
           Name:        "newservice_do_something",
           Description: "Does something useful",
           InputSchema: mcp.InputSchema{
               Type: "object",
               Properties: map[string]mcp.Property{
                   "param": {
                       Type:        "string",
                       Description: "Parameter description",
                   },
               },
               Required: []string{"param"},
           },
       }, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
           param := args["param"].(string)
           return client.DoSomething(ctx)
       })
   }
   ```

4. **Update configuration**

   In `config/base.yaml`:
   ```yaml
   newservice:
     url: ""
     token: ""
     enabled: true
   ```

   In `internal/config/config.go`:
   ```go
   type Config struct {
       // ...
       NewService ServiceConfig `koanf:"newservice"`
   }
   ```

5. **Register in main.go**
   ```go
   import "github.com/axinova-ai/axinova-mcp-server-go/internal/clients/newservice"

   // In main():
   if cfg.NewService.Enabled && cfg.NewService.URL != "" {
       client := newservice.NewClient(
           cfg.NewService.URL,
           cfg.NewService.Token,
           cfg.Timeout.HTTP,
           cfg.TLS.SkipVerify,
       )
       newservice.RegisterTools(mcpServer, client)
       log.Printf("✓ NewService tools registered (%s)", cfg.NewService.URL)
   }
   ```

6. **Update documentation**
   - Add to README.md Available Tools table
   - Add examples to TESTING.md

### Modifying Existing Tools

Tool definitions are in `internal/clients/{service}/tools.go`. Each tool has:
- **Name**: Unique identifier (e.g., `prometheus_query`)
- **Description**: Human-readable explanation
- **InputSchema**: JSON Schema for parameters
- **Handler**: Implementation function

Example modification:
```go
server.RegisterTool(mcp.Tool{
    Name:        "prometheus_query",
    Description: "Execute an instant Prometheus query",
    InputSchema: mcp.InputSchema{
        Type: "object",
        Properties: map[string]mcp.Property{
            "query": {
                Type:        "string",
                Description: "PromQL query expression",
            },
            // Add new parameter
            "timeout": {
                Type:        "string",
                Description: "Query timeout (e.g., '30s')",
            },
        },
        Required: []string{"query"},
    },
}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
    query := args["query"].(string)
    timeout := args["timeout"].(string) // Handle new parameter
    // ... implementation
})
```

### Testing Changes

```bash
# 1. Format and tidy
make fmt
make tidy

# 2. Build
make build

# 3. Test protocol
./test_mcp.sh

# 4. Test specific tool
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}
{"jsonrpc":"2.0","method":"initialized"}
{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"your_tool_name","arguments":{"param":"value"}}}' | ./bin/axinova-mcp-server
```

## Configuration System

### Precedence (lowest to highest)

1. `config/base.yaml` - Defaults
2. `config/{ENV}.yaml` - Environment-specific (dev/prod)
3. Environment variables with `APP_` prefix

### Environment Variable Mapping

```bash
# Flat keys
APP_LOG__LEVEL=debug          → log.level
APP_SERVER__NAME=my-server    → server.name

# Nested with double underscore
APP_PORTAINER__URL=https://...  → portainer.url
APP_TLS__SKIP_VERIFY=true       → tls.skip_verify
```

### TLS Configuration

For internal services with Traefik TLS termination:
```yaml
# config/prod.yaml
tls:
  skip_verify: true
```

Or via environment:
```bash
export APP_TLS__SKIP_VERIFY=true
```

## MCP Protocol Details

### Communication Flow

1. **Initialize** - Client connects and negotiates capabilities
2. **Initialized** - Client confirms ready
3. **Operations** - Client calls tools, lists resources, etc.
4. **Shutdown** - Connection closes

### Available MCP Methods

- `initialize` - Start session
- `initialized` - Confirm initialization
- `tools/list` - Get all available tools
- `tools/call` - Execute a tool
- `resources/list` - List resources (future)
- `resources/read` - Read resource (future)
- `prompts/list` - List prompts (future)
- `prompts/get` - Get prompt (future)

### Tool Call Example

Request:
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "prometheus_query",
    "arguments": {
      "query": "up"
    }
  }
}
```

Response:
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "{\"status\":\"success\",\"data\":{...}}"
      }
    ],
    "isError": false
  }
}
```

## Deployment

### Local Deployment

```bash
# 1. Configure tokens
cp .env.example .env
# Edit .env with your tokens

# 2. Build and run
make build
./bin/axinova-mcp-server
```

### Docker Deployment

```bash
# Build image
make docker-build

# Run with docker-compose
docker-compose up -d

# Check logs
docker-compose logs -f mcp-server
```

### Production Deployment (ax-sas-tools)

See deployment section in main task for SSH deployment to 121.40.188.25.

## Security Considerations

1. **API Tokens**: Never commit tokens to git. Use environment variables or `.env` file (gitignored).

2. **TLS**: Internal services use `tls.skip_verify: true` because TLS is terminated at Traefik.

3. **Tool Safety**: Tools can perform destructive operations (restart containers, delete dashboards). Use with caution.

4. **Network**: Server requires network access to all internal services. Ensure VPN/firewall rules allow connections.

## Troubleshooting

### Common Issues

1. **Build Fails**
   ```bash
   # Clean and retry
   make clean
   go mod tidy
   make build
   ```

2. **TLS Errors in Production**
   ```bash
   # Ensure TLS skip verify is enabled
   export APP_TLS__SKIP_VERIFY=true
   ```

3. **Tool Not Found**
   ```bash
   # List registered tools
   ./test_mcp.sh | grep "name"
   ```

4. **Authentication Errors**
   ```bash
   # Test token manually
   curl -H "Authorization: Bearer $TOKEN" $SERVICE_URL/api/endpoint
   ```

## Code Style

- **Go Formatting**: Use `gofmt` (tabs for indentation)
- **Error Handling**: Always return descriptive errors
- **Logging**: Use `log.Printf` for info, errors go to stderr
- **Comments**: Export functions must have godoc comments
- **Tests**: Place `*_test.go` files next to source

## Resources

- [MCP Specification](https://modelcontextprotocol.io/specification/2025-11-25)
- [MCP GitHub](https://github.com/modelcontextprotocol/modelcontextprotocol)
- [Koanf Documentation](https://github.com/knadh/koanf)
- [Go HTTP Client](https://pkg.go.dev/net/http)

## Conventions

- **Tool Naming**: `{service}_{action}` (e.g., `portainer_list_containers`)
- **Config Keys**: Snake case (e.g., `protocol_version`, `skip_verify`)
- **Environment Variables**: Uppercase with `APP_` prefix
- **File Names**: Snake case (e.g., `client.go`, `tools.go`)
- **Package Names**: Lowercase, no underscores (e.g., `portainer`, `silverbullet`)

## Future Enhancements

- [ ] Add MCP Resources for real-time data streams
- [ ] Implement MCP Prompts for common workflows
- [ ] Add caching for frequently accessed data
- [ ] Support SSE transport in addition to stdio
- [ ] Add metrics and monitoring
- [ ] Implement rate limiting per service
- [ ] Add retry logic with exponential backoff
- [ ] Support batch operations

## Contact

For issues or questions:
- GitHub Issues: https://github.com/axinova-ai/axinova-mcp-server-go/issues
- Internal: Axinova DevOps team
