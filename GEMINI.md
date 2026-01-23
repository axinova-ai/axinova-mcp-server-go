# Axinova MCP Server (Go)

**Service Type:** Backend / Tooling
**Technologies:** Go 1.22+, MCP Protocol (2025-11-25), Koanf, Docker

## Project Overview
An implementation of the Model Context Protocol (MCP) server. It exposes Axinova's internal tools (Portainer, Grafana, Prometheus, SilverBullet, Vikunja) to AI agents via a unified interface.

## Key Directories
- `cmd/server/`: Main entrypoint.
- `internal/mcp/`: MCP protocol implementation (JSON-RPC 2.0).
- `internal/clients/`: Client implementations for external services.
- `internal/config/`: Configuration loading.

## Development Commands
- **Run Locally:** `make run`
- **Build:** `make build`
- **Docker:** `make docker-build`
- **Test:** `make test`

## Configuration
- **Env Vars:** `APP_PORTAINER__URL`, `APP_PORTAINER__TOKEN`, etc.
- **Koanf:** Loads from `config/base.yaml` -> `config/dev.yaml` -> Env.

## Usage
- Can be used with Claude Desktop or other MCP clients.
- Supports `stdio` transport.
