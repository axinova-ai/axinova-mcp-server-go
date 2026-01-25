# Documentation Index

Welcome to the Axinova MCP Server documentation. This index will help you find the right guide for your use case.

## 🚀 Getting Started

Choose your integration method based on your client:

### Native MCP Integration (Recommended)

For clients with built-in MCP protocol support:

- **[Native MCP Integration Guide](NATIVE-MCP-INTEGRATION.md)** - Overview, architecture, and quick start
- **[Claude Desktop Onboarding](onboarding/claude-desktop.md)** - Step-by-step setup for Claude Desktop app
- **[Claude Code Onboarding](onboarding/claude-code.md)** - CLI integration guide
- **[GitHub Copilot Onboarding](onboarding/github-copilot.md)** - VS Code, JetBrains, and CLI setup

**Best for:** Claude Desktop, Claude Code, GitHub Copilot, and any MCP-compliant client

**Benefits:**
- Native protocol support
- Automatic tool discovery
- Real-time streaming
- No API wrappers needed

### API Integration (For Non-MCP Clients)

For platforms without native MCP support:

- **[API Reference](API-REFERENCE.md)** - Complete HTTP API documentation
- **[LLM Integration Guide](LLM-INTEGRATION.md)** - Integration examples for various platforms
- **[Tool Catalog](TOOL-CATALOG.md)** - Complete list of all 38 tools

**Best for:** ChatGPT, Gemini, LangChain, custom applications

**Note:** This uses the HTTP JSON-RPC wrapper, not the native MCP protocol.

## 💡 Examples

Ready-to-use configuration examples:

- **[Claude Desktop Config](examples/claude_desktop_config.json)** - Working `claude_desktop_config.json`
- **[VS Code Settings](examples/vscode_settings.json)** - GitHub Copilot configuration
- **[Installation Script](examples/install-local.sh)** - Automated local installation

## 🗺 Integration Quick Reference

| Client | Transport | Config Location | Guide |
|--------|-----------|----------------|-------|
| **Claude Desktop** | stdio | `~/Library/Application Support/Claude/claude_desktop_config.json` | [Setup](onboarding/claude-desktop.md) |
| **Claude Code** | stdio | Same as Claude Desktop | [Setup](onboarding/claude-code.md) |
| **GitHub Copilot (VS Code)** | stdio | `settings.json` or `.vscode/settings.json` | [Setup](onboarding/github-copilot.md) |
| **GitHub Copilot (JetBrains)** | stdio | IDE Settings → GitHub Copilot → MCP | [Setup](onboarding/github-copilot.md) |
| **ChatGPT / Custom** | HTTP | API endpoint configuration | [API Guide](API-REFERENCE.md) |

## 🛠 Available Services

The MCP server provides tools for these services:

- **Portainer** (8 tools) - Docker container management
- **Grafana** (9 tools) - Monitoring dashboards
- **Prometheus** (7 tools) - Metrics and alerting
- **SilverBullet** (6 tools) - Wiki and knowledge base
- **Vikunja** (8 tools) - Task and project management

**Total: 38 tools**

See [Tool Catalog](TOOL-CATALOG.md) for complete list.

## 🔧 Configuration

### Environment Variables

All services are configured via environment variables with the `APP_` prefix:

```bash
# Portainer
APP_PORTAINER__URL=https://portainer.example.com
APP_PORTAINER__TOKEN=ptr_xxx

# Grafana
APP_GRAFANA__URL=https://grafana.example.com
APP_GRAFANA__TOKEN=glsa_xxx

# Prometheus
APP_PROMETHEUS__URL=https://prometheus.example.com

# SilverBullet
APP_SILVERBULLET__URL=https://notes.example.com
APP_SILVERBULLET__TOKEN=xxx

# Vikunja
APP_VIKUNJA__URL=https://tasks.example.com
APP_VIKUNJA__TOKEN=xxx

# Optional settings
APP_TLS__SKIP_VERIFY=true  # For self-signed certs
APP_LOG__LEVEL=info        # debug, info, warn, error
APP_HTTP__TIMEOUT=30s      # HTTP client timeout
```

**Note:** Use double underscores (`__`) for nested configuration keys.

## 🚀 Installation Options

### Option 1: Download Pre-built Binary

```bash
# macOS
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-macos -o /usr/local/bin/axinova-mcp-server
chmod +x /usr/local/bin/axinova-mcp-server

# Linux
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-linux -o /usr/local/bin/axinova-mcp-server
chmod +x /usr/local/bin/axinova-mcp-server
```

### Option 2: Use Installation Script

```bash
# Run automated installation
curl -fsSL https://raw.githubusercontent.com/axinova-ai/axinova-mcp-server-go/main/docs/examples/install-local.sh | bash
```

### Option 3: Docker

```bash
# Pull image
docker pull ghcr.io/axinova-ai/axinova-mcp-server-go:latest

# Run (stdio mode)
docker run -i --rm \
  -e APP_PORTAINER__TOKEN=ptr_xxx \
  ghcr.io/axinova-ai/axinova-mcp-server-go:latest
```

### Option 4: Build from Source

```bash
# Clone repository
git clone https://github.com/axinova-ai/axinova-mcp-server-go.git
cd axinova-mcp-server-go

# Build
make build

# Install globally
sudo make install
```

## 🔒 Security Best Practices

### Token Management

- **Store tokens securely** - Use environment variables or secure vaults
- **Never commit tokens** - Add config files with tokens to `.gitignore`
- **Use read-only tokens** - Where possible, use viewer/read-only roles
- **Rotate regularly** - Set reminders to rotate tokens every 90 days
- **Separate environments** - Use different tokens for dev/stage/prod

### File Permissions

```bash
# Protect your config files
chmod 600 ~/Library/Application\ Support/Claude/claude_desktop_config.json
chmod 600 .vscode/settings.json
```

## 🐛 Troubleshooting

### Common Issues

**Server not starting:**
- Check binary is executable: `chmod +x /usr/local/bin/axinova-mcp-server`
- Test manually: `echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{...}}' | axinova-mcp-server`

**Tools not appearing in client:**
- Verify config file syntax (valid JSON)
- Check client logs (e.g., `~/Library/Logs/Claude/mcp*.log` for Claude Desktop)
- Restart client completely (quit and relaunch)

**Connection errors to services:**
- Verify service URLs are accessible: `curl -I https://portainer.example.com`
- Check API tokens are valid and not expired
- Ensure `APP_TLS__SKIP_VERIFY=true` for self-signed certs

**Permission errors:**
- Make binary executable: `chmod +x /usr/local/bin/axinova-mcp-server`
- Check ownership: `ls -la /usr/local/bin/axinova-mcp-server`

### Debug Mode

Enable debug logging for troubleshooting:

```bash
APP_LOG__LEVEL=debug
APP_LOG__FORMAT=json  # Optional: structured logs
```

---

## 📚 Operational Documentation

### Deployment & Infrastructure

- **[Deployment Guide](runbooks/DEPLOYMENT.md)** - Production deployment procedures
- **[Testing Guide](runbooks/TESTING.md)** - How to test MCP endpoints
- **[Token Generation](ops/TOKEN_GENERATION_WALKTHROUGH.md)** - Service token setup
- **[Validation Procedures](ops/VALIDATION.md)** - Configuration validation

### Development

- **[CLAUDE.md](CLAUDE.md)** - AI agent development guide
- **[Main README](../README.md)** - Project overview

### Implementation Status

- **[Infrastructure Analysis](ops/INFRASTRUCTURE_ANALYSIS.md)** - Infrastructure review
- **[Issues Fixed](ops/ISSUES_FIXED_AND_REMAINING.md)** - Known issues

---

## 📝 Additional Resources

### External Documentation

- **[MCP Specification](https://spec.modelcontextprotocol.io/)** - Official Model Context Protocol spec
- **[Claude Desktop](https://claude.ai)** - Download Claude Desktop
- **[GitHub Copilot Docs](https://docs.github.com/copilot)** - GitHub Copilot documentation

### Project Links

- **[GitHub Repository](https://github.com/axinova-ai/axinova-mcp-server-go)** - Source code and issues
- **[Production Server](https://mcp.axinova-ai.com)** - Live API endpoint (for HTTP integration)

---

## Navigation

**Quick Links:**
- [← Back to Main README](../README.md)
- [Native MCP Integration Guide →](NATIVE-MCP-INTEGRATION.md)
- [Tool Catalog →](TOOL-CATALOG.md)
- [API Reference →](API-REFERENCE.md)
