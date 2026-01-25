# Native MCP Integration Guide

Complete guide for integrating the Axinova MCP server with MCP-native clients.

## Overview

The Axinova MCP server implements the **Model Context Protocol v2025-11-25** with native support for stdio and SSE transports. This enables seamless integration with:

- ✅ Claude Desktop
- ✅ Claude Code CLI
- ✅ GitHub Copilot (VS Code, JetBrains, CLI)
- ✅ Any MCP-compliant client

## Architecture

The server provides **38 tools across 5 services**:
- **Portainer** (8 tools) - Docker container management
- **Grafana** (9 tools) - Monitoring dashboards
- **Prometheus** (7 tools) - Metrics and alerting
- **SilverBullet** (6 tools) - Wiki and notes
- **Vikunja** (8 tools) - Task management

## Transport Mechanisms

### stdio (Local/Container)
- **Use Case:** Local development, container deployments
- **Clients:** Claude Desktop, Claude Code, GitHub Copilot
- **Configuration:** Command + environment variables
- **Authentication:** Via service tokens in environment

The stdio transport enables the server to run as a subprocess, communicating via standard input/output streams. This is the most common transport for desktop and CLI integrations.

### SSE (Remote/Web)
- **Use Case:** Web-based MCP clients
- **Endpoint:** `https://mcp.axinova-ai.com/api/mcp/v1/sse`
- **Authentication:** Bearer token
- **Status:** Available (requires bridge for Claude Desktop/Code)

Server-Sent Events (SSE) transport enables remote connections over HTTP. Note that Anthropic is deprecating SSE in favor of Streamable HTTP transport.

### HTTP JSON-RPC (Non-MCP Clients)
- **Use Case:** ChatGPT, Gemini, LangChain, custom apps
- **Endpoint:** `https://mcp.axinova-ai.com/api/mcp/v1/call`
- **Note:** This is NOT the native MCP protocol

The HTTP API provides a compatibility layer for platforms that don't support native MCP protocol.

## Quick Start

Choose your client:
- **[Claude Desktop Setup](onboarding/claude-desktop.md)** - macOS and Windows desktop app
- **[Claude Code Setup](onboarding/claude-code.md)** - Command-line interface
- **[GitHub Copilot Setup](onboarding/github-copilot.md)** - VS Code, JetBrains, CLI

## Requirements

### Service Access Tokens

You'll need API tokens for the services you want to access:

- **Portainer:** API token from Settings → API access tokens
- **Grafana:** Service account token from Administration → Service accounts
- **Prometheus:** Usually no auth required for internal deployments
- **SilverBullet:** API token from settings
- **Vikunja:** API token from user settings

### MCP Server Binary or Docker Image

**Option 1: Download Pre-built Binary**
```bash
# macOS
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-macos -o /usr/local/bin/axinova-mcp-server

# Linux
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-linux -o /usr/local/bin/axinova-mcp-server

# Make executable
chmod +x /usr/local/bin/axinova-mcp-server
```

**Option 2: Build from Source**
```bash
cd /path/to/axinova-mcp-server-go
make build
sudo make install
```

**Option 3: Use Docker Image**
```bash
docker pull ghcr.io/axinova-ai/axinova-mcp-server-go:latest
```

### Client Application

Install your preferred MCP client:
- **Claude Desktop:** Download from claude.ai
- **Claude Code:** Install via `npm install -g @anthropic-ai/claude-code`
- **GitHub Copilot:** Install extension in VS Code or JetBrains

## Available Tools

### Portainer Tools (8)
- `portainer_list_containers` - List all containers
- `portainer_get_container` - Get container details
- `portainer_start_container` - Start a stopped container
- `portainer_stop_container` - Stop a running container
- `portainer_restart_container` - Restart a container
- `portainer_get_container_logs` - Retrieve container logs
- `portainer_list_stacks` - List Docker Compose stacks
- `portainer_get_stack` - Get stack details

### Grafana Tools (9)
- `grafana_list_dashboards` - List all dashboards
- `grafana_get_dashboard` - Get dashboard by UID
- `grafana_search_dashboards` - Search dashboards by query
- `grafana_get_dashboard_panels` - Get panels from a dashboard
- `grafana_list_datasources` - List configured data sources
- `grafana_list_folders` - List dashboard folders
- `grafana_list_alerts` - List alert rules
- `grafana_get_alert` - Get alert details
- `grafana_test_datasource` - Test datasource connection

### Prometheus Tools (7)
- `prometheus_query` - Execute instant query
- `prometheus_query_range` - Execute range query
- `prometheus_get_targets` - List scrape targets
- `prometheus_get_alerts` - List active alerts
- `prometheus_get_alert_rules` - List alert rules
- `prometheus_get_metrics` - List available metrics
- `prometheus_get_label_values` - Get values for a label

### SilverBullet Tools (6)
- `silverbullet_list_pages` - List all wiki pages
- `silverbullet_read_page` - Read page content
- `silverbullet_search` - Search across pages
- `silverbullet_get_page_meta` - Get page metadata
- `silverbullet_list_templates` - List page templates
- `silverbullet_query` - Execute dataview query

### Vikunja Tools (8)
- `vikunja_list_projects` - List all projects
- `vikunja_get_project` - Get project details
- `vikunja_list_tasks` - List tasks (filterable)
- `vikunja_get_task` - Get task details
- `vikunja_create_task` - Create new task
- `vikunja_update_task` - Update task
- `vikunja_list_labels` - List available labels
- `vikunja_search_tasks` - Search tasks by text

## Example Usage

Once configured, you can interact with these tools using natural language:

**Container Management:**
- "Show me all running Docker containers"
- "Get the last 100 lines of logs from the grafana container"
- "Restart the postgres container"

**Monitoring:**
- "List all Grafana dashboards"
- "What's the current CPU usage across all services?"
- "Show me active Prometheus alerts"

**Knowledge Management:**
- "Search the wiki for deployment procedures"
- "List all pages in the SilverBullet wiki"

**Task Management:**
- "Create a task to review production logs"
- "List all tasks in the DevOps project"
- "Show me high-priority tasks"

## Security Considerations

### Token Management

- Store tokens in environment variables, never hardcode
- Use separate tokens for dev/stage/prod environments
- Rotate tokens regularly
- Use read-only tokens where possible

### Network Security

- MCP server communicates over HTTPS with internal services
- Set `APP_TLS__SKIP_VERIFY=true` for self-signed certs (dev only)
- Use production-grade certificates in production

### Client Configuration

- Protect your client configuration files (contain tokens)
- On macOS, `claude_desktop_config.json` is in user-only directory
- Consider using environment variables from `.env` files

## Troubleshooting

### Server Not Starting

**Check binary is executable:**
```bash
which axinova-mcp-server
ls -la /usr/local/bin/axinova-mcp-server
```

**Test manually:**
```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | axinova-mcp-server
```

Expected: JSON response with server capabilities

### Tools Not Appearing in Client

**Claude Desktop:**
- Check logs: `~/Library/Logs/Claude/mcp*.log`
- Verify config syntax (valid JSON)
- Restart Claude Desktop completely

**GitHub Copilot:**
- Check VS Code Developer Tools (Help → Toggle Developer Tools)
- Verify Copilot extension is updated
- Check MCP policy is enabled (Enterprise)

### Connection Errors to Services

**Verify service URLs:**
```bash
curl -I https://portainer.axinova-internal.xyz
curl -I https://grafana.axinova-internal.xyz
```

**Check tokens:**
- Tokens must be valid and not expired
- Tokens must have appropriate permissions
- Format: `APP_PORTAINER__TOKEN=ptr_...` (no quotes in env)

### Permission Denied Errors

**Binary not executable:**
```bash
chmod +x /usr/local/bin/axinova-mcp-server
```

**Installation requires sudo:**
```bash
sudo mv axinova-mcp-server /usr/local/bin/
```

## Advanced Configuration

### Custom Service Endpoints

Override default URLs via environment variables:
```bash
APP_PORTAINER__URL=https://custom-portainer.example.com
APP_GRAFANA__URL=https://custom-grafana.example.com
```

### Selective Service Enablement

Disable services by omitting their tokens:
```bash
# Only enable Portainer and Grafana
APP_PORTAINER__TOKEN=ptr_xxx
APP_GRAFANA__TOKEN=glsa_xxx
# Prometheus, SilverBullet, Vikunja will be skipped
```

### Logging Configuration

Adjust log level for debugging:
```bash
APP_LOG__LEVEL=debug  # Default: info
APP_LOG__FORMAT=json  # Default: text
```

### Timeout Configuration

Adjust HTTP client timeouts:
```bash
APP_HTTP__TIMEOUT=30s  # Default: 10s
```

## Migration from HTTP API

If you're currently using the HTTP JSON-RPC API with custom wrappers (e.g., ChatGPT integration), you can migrate to native MCP:

**Before (HTTP API):**
- Custom Python/Node.js wrapper
- Manual tool definitions
- HTTP polling or webhooks

**After (Native MCP):**
- Direct stdio/SSE connection
- Automatic tool discovery
- Built-in streaming support

See individual onboarding guides for migration steps.

## Next Steps

1. **Choose your client** - Select Claude Desktop, Claude Code, or GitHub Copilot
2. **Follow the onboarding guide** - Step-by-step setup instructions
3. **Configure service tokens** - Add API tokens for services you use
4. **Test the integration** - Try example commands
5. **Explore the tools** - Discover what's possible with 38 DevOps tools

## Support

- **Documentation:** This repository's `docs/` directory
- **Issues:** GitHub Issues for bug reports and feature requests
- **MCP Specification:** https://spec.modelcontextprotocol.io/

## Related Documentation

- **[Claude Desktop Setup](onboarding/claude-desktop.md)**
- **[Claude Code Setup](onboarding/claude-code.md)**
- **[GitHub Copilot Setup](onboarding/github-copilot.md)**
- **[Example Configurations](examples/)**
- **[HTTP API Integration](UNIVERSAL-API-INTEGRATION.md)** (for non-MCP clients)
