# Quick Start: Native MCP Integration (Simplest Approach)

**⏱️ Time to setup:** 5 minutes
**💡 Difficulty:** Beginner
**✅ Recommended for:** All MCP clients (Claude Desktop, Claude Code, GitHub Copilot)

---

## Overview

This guide shows the **simplest way** to integrate the Axinova MCP server with MCP-native clients using stdio transport. No HTTP APIs, no wrappers, no complexity.

### What You'll Get

- **38 DevOps tools** available in your IDE/desktop app
- **Natural language interface** to infrastructure
- **Sub-10ms latency** (stdio transport)
- **Zero network overhead**

---

## Prerequisites

Choose ONE of these methods:

### Option A: Download Pre-built Binary (Easiest)

```bash
# macOS
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-macos -o /tmp/axinova-mcp-server
sudo mv /tmp/axinova-mcp-server /usr/local/bin/axinova-mcp-server
sudo chmod +x /usr/local/bin/axinova-mcp-server

# Linux
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-linux -o /tmp/axinova-mcp-server
sudo mv /tmp/axinova-mcp-server /usr/local/bin/axinova-mcp-server
sudo chmod +x /usr/local/bin/axinova-mcp-server
```

### Option B: Build from Source

```bash
cd /path/to/axinova-mcp-server-go
make build
sudo cp bin/axinova-mcp-server /usr/local/bin/
sudo chmod +x /usr/local/bin/axinova-mcp-server
```

### Verify Installation

```bash
which axinova-mcp-server
# Should output: /usr/local/bin/axinova-mcp-server
```

---

## Quick Setup by Client

### Claude Desktop (macOS)

**1. Create/edit config file:**

```bash
code ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

**2. Add this configuration:**

```json
{
  "mcpServers": {
    "axinova": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
        "APP_PORTAINER__TOKEN": "ptr_YOUR_TOKEN",
        "APP_GRAFANA__URL": "https://grafana.axinova-internal.xyz",
        "APP_GRAFANA__TOKEN": "glsa_YOUR_TOKEN",
        "APP_TLS__SKIP_VERIFY": "true",
        "APP_SERVER__API_ENABLED": "false",
        "APP_SERVER__SSE_ENABLED": "false"
      }
    }
  }
}
```

**3. Restart Claude Desktop**

Quit (Cmd+Q) and relaunch.

**4. Test**

Type in Claude: "List all Docker containers"

✅ Done!

---

### Claude Code (CLI)

**1. Add MCP server:**

```bash
claude mcp add axinova --scope user --transport stdio -- /usr/local/bin/axinova-mcp-server
```

**2. Edit config to add environment variables:**

```bash
code ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

Add the `env` section from the Claude Desktop example above.

**3. Test**

```bash
claude code
```

Type: "List all Docker containers"

✅ Done!

---

### GitHub Copilot (VS Code)

**1. Open VS Code settings:**

`Cmd+,` → Search for "copilot mcp"

**2. Edit settings.json:**

```json
{
  "github.copilot.mcp.servers": {
    "axinova": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
        "APP_PORTAINER__TOKEN": "ptr_YOUR_TOKEN",
        "APP_TLS__SKIP_VERIFY": "true",
        "APP_SERVER__API_ENABLED": "false",
        "APP_SERVER__SSE_ENABLED": "false"
      }
    }
  }
}
```

**3. Restart VS Code**

**4. Test**

Open Copilot Chat (Cmd+I) and type: "List Docker containers using Portainer"

✅ Done!

---

## Environment Variables Reference

### Required

Only configure the services you want to use:

```bash
# Portainer (Docker management)
APP_PORTAINER__URL=https://portainer.axinova-internal.xyz
APP_PORTAINER__TOKEN=ptr_xxx

# Grafana (Dashboards)
APP_GRAFANA__URL=https://grafana.axinova-internal.xyz
APP_GRAFANA__TOKEN=glsa_xxx

# Prometheus (Metrics)
APP_PROMETHEUS__URL=https://prometheus.axinova-internal.xyz

# SilverBullet (Wiki)
APP_SILVERBULLET__URL=https://silverbullet.axinova-internal.xyz
APP_SILVERBULLET__TOKEN=xxx

# Vikunja (Tasks)
APP_VIKUNJA__URL=https://vikunja.axinova-internal.xyz
APP_VIKUNJA__TOKEN=xxx
```

### Optional (Recommended)

```bash
# For self-signed certificates (internal services)
APP_TLS__SKIP_VERIFY=true

# Disable unused transports (stdio-only mode)
APP_SERVER__API_ENABLED=false   # Disable HTTP API
APP_SERVER__SSE_ENABLED=false   # Disable SSE transport
APP_SERVER__HTTP_ENABLED=false  # Disable health server
```

---

## Testing

### Quick Tests

Try these commands in your client:

**Portainer:**
- "Show all running containers"
- "Get logs for the grafana container"

**Grafana:**
- "List all dashboards"

**Prometheus:**
- "What's the CPU usage?"

**SilverBullet:**
- "List all wiki pages"

**Vikunja:**
- "Show all projects"

### Expected Behavior

✅ Client discovers 38 tools automatically
✅ Tools execute in < 100ms
✅ No HTTP requests (stdio = local)
✅ No authentication errors

---

## Troubleshooting

### Tools not appearing?

**Check logs:**
```bash
# Claude Desktop
tail -f ~/Library/Logs/Claude/mcp*.log
```

**Common fixes:**
1. Verify config is valid JSON (use JSONLint)
2. Restart client completely (Quit, not just close)
3. Check binary exists: `which axinova-mcp-server`

### Connection errors?

**Fix:** Add to config:
```json
"APP_TLS__SKIP_VERIFY": "true"
```

### Permission denied?

**Fix:**
```bash
sudo chmod +x /usr/local/bin/axinova-mcp-server
```

---

## Why This Approach is Better

### vs HTTP API Integration

| Aspect | stdio (Native) | HTTP API |
|--------|---------------|----------|
| **Speed** | 10-20ms | 100-300ms |
| **Setup** | 1 config file | TypeScript + npm + build |
| **Security** | Local process | Network + tokens |
| **Complexity** | ⭐ Simple | ⭐⭐⭐⭐ Complex |

### vs .claude-plugin TypeScript Wrapper

| Feature | Native MCP | .claude-plugin |
|---------|-----------|----------------|
| **Performance** | 10-20ms | 100-300ms |
| **Dependencies** | None | Node.js + npm |
| **Maintenance** | None | Update TypeScript defs |
| **Auto-discovery** | ✅ Yes | ❌ Manual registration |

**Native MCP is 10-30x faster and 90% simpler.**

---

## Next Steps

### Learn More

- **Full guide:** [Native MCP Integration](NATIVE-MCP-INTEGRATION.md)
- **Detailed setup:** [Claude Desktop Onboarding](onboarding/claude-desktop.md)
- **All tools:** [Tool Catalog](TOOL-CATALOG.md)

### Advanced Configuration

- **Multiple environments:** Configure dev/prod separately
- **Selective services:** Only enable what you need
- **Custom timeouts:** Adjust HTTP timeouts
- **Debug logging:** Set `APP_LOG__LEVEL=debug`

---

## Summary

**You just learned the simplest way to integrate MCP:**

1. ✅ Install binary (1 command)
2. ✅ Add config (1 file)
3. ✅ Restart client (1 click)
4. ✅ Use 38 tools (natural language)

**Total time:** 5 minutes
**Complexity:** Minimal
**Performance:** Optimal

🎉 **You're done!** Start asking questions about your infrastructure in natural language.

---

**Questions?** See [Troubleshooting](onboarding/claude-desktop.md#troubleshooting) or [open an issue](https://github.com/axinova-ai/axinova-mcp-server-go/issues).
