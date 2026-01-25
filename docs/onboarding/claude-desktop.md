# Claude Desktop Integration

Step-by-step guide for adding the Axinova MCP server to Claude Desktop.

## Prerequisites

- **Claude Desktop** installed (download from claude.ai)
- **Access to internal services** (Portainer, Grafana, etc.)
- **Service API tokens** for the services you want to use
- **MCP server binary** or Docker image

## Overview

Claude Desktop supports the Model Context Protocol (MCP) via stdio transport. The MCP server runs as a subprocess that Claude Desktop manages automatically - starting it when Claude launches and stopping it on exit.

**Key Features:**
- Automatic server lifecycle management
- 38 DevOps tools available in conversations
- Natural language interface to infrastructure
- Persistent configuration across restarts

## Installation Methods

### Option 1: Local Binary Installation (Recommended)

Best for: Development and personal use

#### Step 1: Install the MCP Server

**macOS:**
```bash
# Download latest release
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-macos -o /tmp/axinova-mcp-server

# Install to system path
sudo mv /tmp/axinova-mcp-server /usr/local/bin/axinova-mcp-server
sudo chmod +x /usr/local/bin/axinova-mcp-server

# Verify installation
/usr/local/bin/axinova-mcp-server --help
```

**Or build from source:**
```bash
cd /path/to/axinova-mcp-server-go
make build
sudo make install
```

#### Step 2: Locate Claude Desktop Config File

**macOS:**
```bash
~/Library/Application Support/Claude/claude_desktop_config.json
```

**Windows:**
```bash
%APPDATA%\Claude\claude_desktop_config.json
```

If the file doesn't exist, create it with `{}` as initial content.

#### Step 3: Configure Claude Desktop

Edit the config file and add the MCP server:

```json
{
  "mcpServers": {
    "axinova-tools": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "ENV": "prod",
        "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
        "APP_PORTAINER__TOKEN": "ptr_YOUR_TOKEN_HERE",
        "APP_GRAFANA__URL": "https://grafana.axinova-internal.xyz",
        "APP_GRAFANA__TOKEN": "glsa_YOUR_TOKEN_HERE",
        "APP_PROMETHEUS__URL": "https://prometheus.axinova-internal.xyz",
        "APP_SILVERBULLET__URL": "https://notes.axinova-internal.xyz",
        "APP_SILVERBULLET__TOKEN": "YOUR_TOKEN_HERE",
        "APP_VIKUNJA__URL": "https://tasks.axinova-internal.xyz",
        "APP_VIKUNJA__TOKEN": "YOUR_TOKEN_HERE",
        "APP_TLS__SKIP_VERIFY": "true"
      }
    }
  }
}
```

**Configuration Notes:**

- **`command`**: Absolute path to the MCP server binary
- **`env`**: Environment variables for configuration
  - Use `__` (double underscore) for nested config keys
  - Example: `APP_PORTAINER__TOKEN` sets `portainer.token`
- **Token format:**
  - Portainer: `ptr_...`
  - Grafana: `glsa_...` (service account token)
  - SilverBullet: Token from settings
  - Vikunja: Token from user settings
- **TLS Skip Verify:** Set to `"true"` for self-signed certificates (dev/internal)

#### Step 4: Restart Claude Desktop

**Important:** Completely quit and restart Claude Desktop:

**macOS:**
1. Cmd+Q to quit (or Claude → Quit Claude)
2. Relaunch from Applications

**Windows:**
1. Right-click system tray icon → Exit
2. Relaunch from Start menu

The MCP server will start automatically when Claude Desktop launches.

#### Step 5: Verify Connection

1. Open Claude Desktop
2. Start a new conversation
3. Type: **"List all Docker containers"**

Claude should recognize the `portainer_list_containers` tool and execute it, showing you the list of containers.

If tools appear, congratulations! The integration is working.

### Option 2: Docker Container

Best for: Consistent environment, easy updates

#### Configuration

```json
{
  "mcpServers": {
    "axinova-tools": {
      "command": "docker",
      "args": [
        "run",
        "-i",
        "--rm",
        "--network", "host",
        "-e", "APP_PORTAINER__URL=https://portainer.axinova-internal.xyz",
        "-e", "APP_PORTAINER__TOKEN=ptr_YOUR_TOKEN_HERE",
        "-e", "APP_GRAFANA__URL=https://grafana.axinova-internal.xyz",
        "-e", "APP_GRAFANA__TOKEN=glsa_YOUR_TOKEN_HERE",
        "-e", "APP_PROMETHEUS__URL=https://prometheus.axinova-internal.xyz",
        "-e", "APP_TLS__SKIP_VERIFY=true",
        "ghcr.io/axinova-ai/axinova-mcp-server-go:latest"
      ]
    }
  }
}
```

**Docker Notes:**
- `-i`: Interactive mode (required for stdio)
- `--rm`: Remove container after exit
- `--network host`: Access to localhost services (optional)
- Each environment variable needs separate `-e` flag

#### Restart and Verify

Same as Option 1 - quit and restart Claude Desktop, then test with "List all Docker containers".

## Configuration Examples

### Minimal Configuration (Portainer Only)

```json
{
  "mcpServers": {
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

Only Portainer tools will be available. Other services are automatically disabled if tokens are missing.

### Full Configuration (All Services)

```json
{
  "mcpServers": {
    "axinova-tools": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "ENV": "prod",
        "APP_LOG__LEVEL": "info",
        "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
        "APP_PORTAINER__TOKEN": "ptr_YOUR_TOKEN_HERE",
        "APP_GRAFANA__URL": "https://grafana.axinova-internal.xyz",
        "APP_GRAFANA__TOKEN": "glsa_YOUR_TOKEN_HERE",
        "APP_PROMETHEUS__URL": "https://prometheus.axinova-internal.xyz",
        "APP_SILVERBULLET__URL": "https://notes.axinova-internal.xyz",
        "APP_SILVERBULLET__TOKEN": "YOUR_TOKEN_HERE",
        "APP_VIKUNJA__URL": "https://tasks.axinova-internal.xyz",
        "APP_VIKUNJA__TOKEN": "YOUR_TOKEN_HERE",
        "APP_TLS__SKIP_VERIFY": "true",
        "APP_HTTP__TIMEOUT": "30s"
      }
    }
  }
}
```

All 38 tools across 5 services will be available.

### Multiple MCP Servers

You can configure multiple MCP servers in Claude Desktop:

```json
{
  "mcpServers": {
    "axinova-tools": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "APP_PORTAINER__TOKEN": "ptr_YOUR_TOKEN_HERE"
      }
    },
    "other-mcp-server": {
      "command": "/usr/local/bin/other-server",
      "env": {}
    }
  }
}
```

## Testing the Integration

### Basic Tests

**Portainer:**
```
You: Show me all Docker containers
```
Expected: List of containers with status, names, and IDs

```
You: Get the last 50 lines of logs from the grafana container
```
Expected: Recent log entries from the container

```
You: Restart the postgres container
```
Expected: Confirmation of restart

**Grafana:**
```
You: List all Grafana dashboards
```
Expected: List of dashboards with titles and UIDs

```
You: What's on the "Container Resources" dashboard?
```
Expected: Dashboard panels and their queries

**Prometheus:**
```
You: What's the current CPU usage across all services?
```
Expected: CPU metrics from Prometheus

```
You: Show me memory usage for the last hour
```
Expected: Memory usage graph/data

**SilverBullet:**
```
You: List all wiki pages
```
Expected: List of page names

```
You: Search the wiki for "deployment"
```
Expected: Pages containing "deployment"

**Vikunja:**
```
You: Show me all projects
```
Expected: List of project names

```
You: Create a task "Review production logs" in the DevOps project
```
Expected: Confirmation of task creation

### Advanced Tests

**Multi-step workflows:**
```
You: Find all containers that are unhealthy, get their logs, and create tasks in Vikunja to investigate each one
```

Claude should orchestrate multiple tool calls across Portainer and Vikunja.

**Cross-service queries:**
```
You: Check if there are any Prometheus alerts firing, then look up the related Grafana dashboard
```

Claude should query Prometheus for alerts, then search Grafana for relevant dashboards.

## Troubleshooting

### Tools Not Appearing

**Symptom:** Claude doesn't recognize MCP tools in conversations

**Solutions:**

1. **Check Claude Desktop logs:**
   ```bash
   # macOS
   tail -f ~/Library/Logs/Claude/mcp*.log
   ```

2. **Verify config file syntax:**
   - Must be valid JSON (use JSONLint.com to validate)
   - Check for missing commas, quotes, brackets
   - Environment variable values should be strings (quoted)

3. **Verify binary path:**
   ```bash
   which axinova-mcp-server
   # Should output: /usr/local/bin/axinova-mcp-server

   ls -la /usr/local/bin/axinova-mcp-server
   # Should show executable permissions: -rwxr-xr-x
   ```

4. **Test server manually:**
   ```bash
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | /usr/local/bin/axinova-mcp-server
   ```

   Expected: JSON response with server capabilities

5. **Completely restart Claude Desktop:**
   - Quit (Cmd+Q), don't just close windows
   - Wait 5 seconds
   - Relaunch

### Connection Errors to Services

**Symptom:** Tools fail with connection errors

**Solutions:**

1. **Verify service URLs are accessible:**
   ```bash
   curl -I https://portainer.axinova-internal.xyz
   curl -I https://grafana.axinova-internal.xyz
   ```

2. **Check API tokens are valid:**
   - Portainer: Settings → API access tokens (verify token exists)
   - Grafana: Administration → Service accounts (verify token hasn't expired)

3. **Test with curl:**
   ```bash
   # Portainer
   curl -H "X-API-Key: ptr_YOUR_TOKEN" https://portainer.axinova-internal.xyz/api/endpoints

   # Grafana
   curl -H "Authorization: Bearer glsa_YOUR_TOKEN" https://grafana.axinova-internal.xyz/api/dashboards
   ```

4. **Check TLS settings:**
   - If using self-signed certificates, ensure `APP_TLS__SKIP_VERIFY` is set to `"true"`

### Environment Variables Not Working

**Symptom:** Server starts but can't connect to services

**Common mistakes:**

1. **Missing quotes around values:**
   ```json
   // Wrong:
   "APP_PORTAINER__TOKEN": ptr_xxx

   // Correct:
   "APP_PORTAINER__TOKEN": "ptr_xxx"
   ```

2. **Single underscore instead of double:**
   ```json
   // Wrong:
   "APP_PORTAINER_TOKEN": "ptr_xxx"

   // Correct:
   "APP_PORTAINER__TOKEN": "ptr_xxx"
   ```

3. **Environment variables in wrong place:**
   ```json
   // Wrong:
   {
     "mcpServers": {
       "axinova-tools": {
         "command": "/usr/local/bin/axinova-mcp-server",
         "args": ["APP_PORTAINER__TOKEN=ptr_xxx"]  // Wrong!
       }
     }
   }

   // Correct:
   {
     "mcpServers": {
       "axinova-tools": {
         "command": "/usr/local/bin/axinova-mcp-server",
         "env": {
           "APP_PORTAINER__TOKEN": "ptr_xxx"  // Correct!
         }
       }
     }
   }
   ```

### Permission Denied Errors

**Symptom:** Server fails to start with permission errors

**Solutions:**

1. **Make binary executable:**
   ```bash
   chmod +x /usr/local/bin/axinova-mcp-server
   ```

2. **Check file ownership:**
   ```bash
   ls -la /usr/local/bin/axinova-mcp-server
   # Should be owned by root or your user
   ```

3. **Reinstall with proper permissions:**
   ```bash
   sudo mv axinova-mcp-server /usr/local/bin/
   sudo chmod 755 /usr/local/bin/axinova-mcp-server
   ```

### Server Crashes or Exits Immediately

**Symptom:** Server starts but immediately exits

**Solutions:**

1. **Check logs for errors:**
   ```bash
   tail -100 ~/Library/Logs/Claude/mcp*.log
   ```

2. **Test with minimal config:**
   - Remove all environment variables except one service
   - Test with just Portainer

3. **Verify server version:**
   ```bash
   /usr/local/bin/axinova-mcp-server --version
   ```

4. **Try running server directly:**
   ```bash
   APP_PORTAINER__URL=https://portainer.example.com \
   APP_PORTAINER__TOKEN=ptr_xxx \
   /usr/local/bin/axinova-mcp-server
   ```

   Then send an initialize message via stdin to see error output.

## Security Best Practices

### Protect Your Config File

The config file contains sensitive API tokens:

```bash
# macOS - restrict to your user
chmod 600 ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

### Use Read-Only Tokens Where Possible

- Grafana: Create viewer-only service accounts
- Portainer: Use read-only API tokens if you only need monitoring

### Separate Tokens Per Environment

Don't use production tokens for development:

```json
{
  "mcpServers": {
    "axinova-dev": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "ENV": "dev",
        "APP_PORTAINER__TOKEN": "ptr_dev_token"
      }
    },
    "axinova-prod": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "ENV": "prod",
        "APP_PORTAINER__TOKEN": "ptr_prod_token"
      }
    }
  }
}
```

Then enable/disable as needed by commenting out configs.

### Rotate Tokens Regularly

Set reminders to rotate API tokens every 90 days.

## Updating the MCP Server

### Binary Installation

```bash
# Download new version
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-macos -o /tmp/axinova-mcp-server

# Replace existing binary
sudo mv /tmp/axinova-mcp-server /usr/local/bin/axinova-mcp-server
sudo chmod +x /usr/local/bin/axinova-mcp-server

# Restart Claude Desktop
```

### Docker Installation

Docker will automatically pull the latest image on next launch if using `:latest` tag.

To update manually:
```bash
docker pull ghcr.io/axinova-ai/axinova-mcp-server-go:latest
```

Then restart Claude Desktop.

## Uninstalling

### Remove MCP Server from Claude Desktop

Edit `~/Library/Application Support/Claude/claude_desktop_config.json` and remove the `axinova-tools` entry:

```json
{
  "mcpServers": {
    // Remove this entire block:
    // "axinova-tools": { ... }
  }
}
```

Restart Claude Desktop.

### Remove Binary

```bash
sudo rm /usr/local/bin/axinova-mcp-server
```

### Remove Docker Image

```bash
docker rmi ghcr.io/axinova-ai/axinova-mcp-server-go:latest
```

## Next Steps

- **Explore more tools:** See [NATIVE-MCP-INTEGRATION.md](../NATIVE-MCP-INTEGRATION.md) for complete tool list
- **Try other clients:** [Claude Code](claude-code.md) or [GitHub Copilot](github-copilot.md)
- **Read examples:** See [examples/](../examples/) for more configuration patterns

## Support

- **GitHub Issues:** Report bugs or request features
- **Documentation:** Full docs in this repository
- **MCP Spec:** https://spec.modelcontextprotocol.io/
