# Claude Code Integration

Add the Axinova MCP server to Claude Code CLI using native MCP protocol.

## Prerequisites

- **Claude Code CLI** installed (`npm install -g @anthropic-ai/claude-code`)
- **Service API tokens** for the services you want to use
- **MCP server binary** or Docker image

## Overview

Claude Code is a command-line interface for Claude that supports MCP natively via stdio transport. MCP servers are configured in the same config file as Claude Desktop, making it easy to share configurations.

**Key Features:**
- Same configuration as Claude Desktop
- Interactive CLI for adding MCP servers
- User and project-level scopes
- Access to 38 DevOps tools from the command line

## Installation

### Step 1: Install MCP Server Binary

**macOS:**
```bash
# Download latest release
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-macos -o /usr/local/bin/axinova-mcp-server
chmod +x /usr/local/bin/axinova-mcp-server
```

**Linux:**
```bash
# Download latest release
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-linux -o /usr/local/bin/axinova-mcp-server
chmod +x /usr/local/bin/axinova-mcp-server
```

**Or build from source:**
```bash
cd /path/to/axinova-mcp-server-go
make build
sudo make install
```

Verify installation:
```bash
which axinova-mcp-server
# Expected: /usr/local/bin/axinova-mcp-server
```

### Step 2: Add MCP Server to Claude Code

You have three methods to configure the MCP server:

## Configuration Methods

### Method 1: Interactive CLI (Recommended)

The simplest way to add an MCP server:

```bash
claude mcp add axinova-tools --scope user --transport stdio -- \
  /usr/local/bin/axinova-mcp-server
```

**Explanation:**
- `axinova-tools` - Name for this MCP server
- `--scope user` - Available in all projects (see Scope Options below)
- `--transport stdio` - Use stdio transport (default)
- `--` - Separator before command
- `/usr/local/bin/axinova-mcp-server` - Path to binary

This creates the configuration automatically. However, this method doesn't allow setting environment variables directly, so you'll need to edit the config file afterward to add service tokens.

**After running the command:**

1. Open the config file:
   ```bash
   # macOS
   code ~/Library/Application\ Support/Claude/claude_desktop_config.json

   # Or use any editor
   vim ~/Library/Application\ Support/Claude/claude_desktop_config.json
   ```

2. Add environment variables:
   ```json
   {
     "mcpServers": {
       "axinova-tools": {
         "command": "/usr/local/bin/axinova-mcp-server",
         "env": {
           "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
           "APP_PORTAINER__TOKEN": "ptr_YOUR_TOKEN_HERE",
           "APP_GRAFANA__URL": "https://grafana.axinova-internal.xyz",
           "APP_GRAFANA__TOKEN": "glsa_YOUR_TOKEN_HERE"
         }
       }
     }
   }
   ```

### Method 2: JSON Configuration (Most Flexible)

For complete control over environment variables:

```bash
claude mcp add-json '{
  "axinova-tools": {
    "type": "stdio",
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
}'
```

**Tips:**
- Use single quotes around the JSON to avoid shell escaping issues
- Make sure JSON is valid (use JSONLint if needed)
- This method sets everything in one command

### Method 3: Manual File Edit

Directly edit the config file:

**Location:**
- **macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux:** `~/.config/Claude/claude_desktop_config.json`

**Content:**
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

This is the same format as Claude Desktop, so you can share configurations between both.

## Scope Options

MCP servers can be configured at different scopes:

### User Scope (Recommended)

Available in all projects and directories:

```bash
claude mcp add axinova-tools --scope user -- /usr/local/bin/axinova-mcp-server
```

**Pros:**
- Always available, no matter where you run `claude code`
- One configuration for all projects
- Shared with Claude Desktop

**Cons:**
- Same configuration everywhere (can't customize per project)

### Project Scope

Only available in the current directory:

```bash
cd /path/to/your/project
claude mcp add axinova-tools --scope project -- /usr/local/bin/axinova-mcp-server
```

This creates `.claude/mcp_config.json` in the current directory.

**Pros:**
- Project-specific tokens and URLs
- Can override user-level config
- Can commit to git (if tokens are in env vars, not hardcoded)

**Cons:**
- Only works in this directory
- Need to configure for each project

## Verification

### List Configured Servers

```bash
claude mcp list
```

**Expected output:**
```
MCP Servers:
  axinova-tools (stdio) - /usr/local/bin/axinova-mcp-server
```

### Test in Claude Code

```bash
claude code
```

In the Claude Code interactive session:

```
You: List all Docker containers
```

Claude should invoke the `portainer_list_containers` tool and show you the containers.

If tools work, the integration is successful!

## Usage Examples

Once configured, you can use natural language to interact with your infrastructure:

### Container Management

```bash
claude code
```

```
You: Show running containers and their status

You: Get the last 100 lines of logs from the grafana container

You: What containers are consuming the most CPU?

You: Restart the postgres container
```

### Monitoring and Observability

```
You: List all Grafana dashboards

You: What's on the "Container Resources" dashboard?

You: Show me Prometheus alerts that are firing

You: What's the current memory usage across all services?
```

### Knowledge Base

```
You: List all wiki pages in SilverBullet

You: Search the wiki for "deployment procedures"

You: What documentation do we have about database backups?
```

### Task Management

```
You: Show me all projects in Vikunja

You: List high-priority tasks

You: Create a task "Review production logs" in the DevOps project

You: What tasks are due this week?
```

### Multi-Tool Workflows

```
You: Check for any failing containers, get their logs, and create a task to investigate

You: Find all Prometheus alerts, then look up related Grafana dashboards

You: Search the wiki for runbook about the current alerts
```

Claude Code will orchestrate multiple tool calls to complete complex requests.

## Configuration Examples

### Minimal (Portainer Only)

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

### Docker-based Installation

```json
{
  "mcpServers": {
    "axinova-tools": {
      "command": "docker",
      "args": [
        "run", "-i", "--rm",
        "-e", "APP_PORTAINER__URL=https://portainer.axinova-internal.xyz",
        "-e", "APP_PORTAINER__TOKEN=ptr_YOUR_TOKEN_HERE",
        "-e", "APP_TLS__SKIP_VERIFY=true",
        "ghcr.io/axinova-ai/axinova-mcp-server-go:latest"
      ]
    }
  }
}
```

### Environment-Specific Configurations

You can configure multiple servers for different environments:

```json
{
  "mcpServers": {
    "axinova-dev": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "ENV": "dev",
        "APP_PORTAINER__URL": "https://portainer-dev.example.com",
        "APP_PORTAINER__TOKEN": "ptr_dev_token"
      }
    },
    "axinova-prod": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "ENV": "prod",
        "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
        "APP_PORTAINER__TOKEN": "ptr_prod_token"
      }
    }
  }
}
```

Then specify which to use in conversation:
```
You: Using axinova-prod, list all containers
```

## Managing MCP Servers

### List Servers

```bash
claude mcp list
```

### Remove a Server

```bash
claude mcp remove axinova-tools
```

This removes the server from the configuration.

### Update Server Configuration

Claude Code doesn't have a direct "update" command, so:

1. Remove the old configuration:
   ```bash
   claude mcp remove axinova-tools
   ```

2. Add with new configuration:
   ```bash
   claude mcp add-json '{...new config...}'
   ```

Or manually edit the config file.

## Troubleshooting

### Server Not Listed

**Problem:** `claude mcp list` doesn't show the server

**Solution:**
1. Check config file exists and is valid JSON:
   ```bash
   # macOS
   cat ~/Library/Application\ Support/Claude/claude_desktop_config.json | jq .

   # Linux
   cat ~/.config/Claude/claude_desktop_config.json | jq .
   ```

2. Verify `mcpServers` key exists with your server entry

### Tools Not Available in Conversation

**Problem:** Claude doesn't recognize MCP tools

**Solution:**
1. **Restart Claude Code session:**
   - Exit current session (Ctrl+C or type "exit")
   - Start new session: `claude code`

2. **Test server manually:**
   ```bash
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | /usr/local/bin/axinova-mcp-server
   ```

   Should return JSON response with server capabilities.

3. **Check binary path:**
   ```bash
   which axinova-mcp-server
   ls -la /usr/local/bin/axinova-mcp-server
   ```

4. **Verify environment variables:**
   - Check for typos in token names
   - Ensure values are quoted strings in JSON
   - Use `__` (double underscore) for nested keys

### Connection Errors to Services

**Problem:** Tools fail with HTTP errors or timeouts

**Solution:**
1. **Test service URLs:**
   ```bash
   curl -I https://portainer.axinova-internal.xyz
   curl -I https://grafana.axinova-internal.xyz
   ```

2. **Verify tokens are valid:**
   ```bash
   # Portainer
   curl -H "X-API-Key: ptr_YOUR_TOKEN" \
     https://portainer.axinova-internal.xyz/api/endpoints

   # Grafana
   curl -H "Authorization: Bearer glsa_YOUR_TOKEN" \
     https://grafana.axinova-internal.xyz/api/dashboards
   ```

3. **Check TLS settings:**
   - For self-signed certs: `"APP_TLS__SKIP_VERIFY": "true"`

### Permission Denied

**Problem:** Server fails to execute

**Solution:**
```bash
# Make binary executable
chmod +x /usr/local/bin/axinova-mcp-server

# Verify permissions
ls -la /usr/local/bin/axinova-mcp-server
# Should show: -rwxr-xr-x
```

### JSON Parse Errors

**Problem:** `claude mcp add-json` fails with JSON error

**Solution:**
1. **Validate JSON:**
   - Use JSONLint.com
   - Check for missing commas, quotes, brackets

2. **Escape shell characters:**
   ```bash
   # Use single quotes for outer JSON
   claude mcp add-json '{"server": {"command": "/path"}}'

   # Not double quotes (shell will interpret)
   ```

3. **Save to file and read:**
   ```bash
   cat > /tmp/mcp-config.json <<'EOF'
   {
     "axinova-tools": {
       "command": "/usr/local/bin/axinova-mcp-server",
       "env": {...}
     }
   }
   EOF

   claude mcp add-json "$(cat /tmp/mcp-config.json)"
   ```

## Advanced Usage

### Using with CI/CD

Claude Code can be used in scripts and CI/CD pipelines:

```bash
#!/bin/bash
# deploy-check.sh

# Configure MCP server for CI
export CLAUDE_CONFIG=$(cat <<EOF
{
  "mcpServers": {
    "axinova-tools": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "APP_PORTAINER__TOKEN": "$PORTAINER_TOKEN"
      }
    }
  }
}
EOF
)

# Ask Claude to check deployment
claude code --prompt "Check if all containers in production are healthy and running"
```

### Project-Specific Tokens

For team projects, store tokens in `.env` files (not committed to git):

**.env:**
```bash
PORTAINER_TOKEN=ptr_xxx
GRAFANA_TOKEN=glsa_xxx
```

**.claude/mcp_config.json:**
```json
{
  "mcpServers": {
    "axinova-tools": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "APP_PORTAINER__TOKEN": "${PORTAINER_TOKEN}",
        "APP_GRAFANA__TOKEN": "${GRAFANA_TOKEN}"
      }
    }
  }
}
```

Then source `.env` before running Claude Code:
```bash
source .env
claude code
```

**Note:** Check if Claude Code supports env var expansion in config files. If not, use a wrapper script to render the config.

### Custom Logging

Debug MCP communication:

```json
{
  "mcpServers": {
    "axinova-tools": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "APP_LOG__LEVEL": "debug",
        "APP_LOG__FORMAT": "json"
      }
    }
  }
}
```

Logs will include detailed MCP protocol messages.

## Best Practices

### 1. Use User Scope for Personal Work

```bash
claude mcp add axinova-tools --scope user -- /usr/local/bin/axinova-mcp-server
```

Then edit config file once with your personal tokens.

### 2. Use Project Scope for Team Work

```bash
cd /your/team/project
claude mcp add axinova-tools --scope project -- /usr/local/bin/axinova-mcp-server
```

Commit `.claude/mcp_config.json` to git without tokens, document how to add tokens in README.

### 3. Separate Dev and Prod

```json
{
  "mcpServers": {
    "axinova-dev": {...},
    "axinova-prod": {...}
  }
}
```

Explicitly specify which environment in your prompts.

### 4. Protect Your Config

```bash
# Restrict config file permissions
chmod 600 ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

### 5. Version Control Friendly

For team projects, use a template:

**mcp_config.json.template:**
```json
{
  "mcpServers": {
    "axinova-tools": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "APP_PORTAINER__TOKEN": "REPLACE_WITH_YOUR_TOKEN"
      }
    }
  }
}
```

Commit the template, not the actual config with tokens.

## Next Steps

- **Explore Claude Code features:** `claude --help`
- **Try other clients:** [Claude Desktop](claude-desktop.md) or [GitHub Copilot](github-copilot.md)
- **Learn more about MCP:** [NATIVE-MCP-INTEGRATION.md](../NATIVE-MCP-INTEGRATION.md)
- **See more examples:** [examples/](../examples/)

## Support

- **GitHub Issues:** Report bugs or request features
- **Documentation:** Full docs in this repository
- **MCP Spec:** https://spec.modelcontextprotocol.io/
