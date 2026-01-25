# GitHub Copilot Integration

Enable the Axinova MCP server in GitHub Copilot (VS Code, JetBrains, CLI).

## Prerequisites

- **GitHub Copilot subscription** (Free, Pro, or Business/Enterprise)
- **VS Code 1.102+** (or JetBrains IDE with Copilot plugin)
- **MCP server binary** or Docker image
- **For Business/Enterprise:** "MCP servers in Copilot" policy enabled by org admin

## Overview

GitHub Copilot added native MCP support in 2026, allowing you to extend Copilot with custom tools via the Model Context Protocol. The Axinova MCP server provides 38 DevOps tools directly in your IDE.

**Key Features:**
- MCP tools available in Copilot Chat
- Natural language interface to infrastructure
- IDE-integrated (VS Code, JetBrains)
- CLI support via `gh copilot`

**Important:** MCP support in GitHub Copilot is currently in preview and may have limitations.

## VS Code Setup

### Step 1: Install MCP Server Binary

**macOS:**
```bash
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-macos -o /usr/local/bin/axinova-mcp-server
chmod +x /usr/local/bin/axinova-mcp-server
```

**Linux:**
```bash
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-linux -o /usr/local/bin/axinova-mcp-server
chmod +x /usr/local/bin/axinova-mcp-server
```

**Windows:**
```powershell
# Download from GitHub releases
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-windows.exe -o C:\Program Files\axinova-mcp-server.exe
```

Verify installation:
```bash
which axinova-mcp-server
# Expected: /usr/local/bin/axinova-mcp-server
```

### Step 2: Configure VS Code Settings

You have two options: **User Settings** (global) or **Workspace Settings** (project-specific).

#### Option A: User Settings (Global)

Best for: Personal use, all projects

1. **Open Settings:**
   - macOS: `Cmd+,`
   - Windows/Linux: `Ctrl+,`

2. **Search for:** `copilot mcp`

3. **Edit settings.json:**
   - Click "Edit in settings.json" link
   - Or open: `~/Library/Application Support/Code/User/settings.json` (macOS)

4. **Add MCP server configuration:**

```json
{
  "github.copilot.mcp.servers": {
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

#### Option B: Workspace Settings (Project-Specific)

Best for: Team projects, environment-specific configs

1. **Create `.vscode/settings.json` in your project:**

```bash
mkdir -p /your/project/.vscode
```

2. **Add configuration:**

```json
{
  "github.copilot.mcp.servers": {
    "axinova-tools": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "ENV": "prod",
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

**Important:** Don't commit tokens to git! Use placeholders and document how team members should add their tokens.

### Step 3: Restart VS Code

Close and reopen VS Code for changes to take effect.

### Step 4: Verify Integration

1. **Open Copilot Chat:**
   - macOS: `Cmd+I`
   - Windows/Linux: `Ctrl+I`

2. **Test with a command:**
   ```
   List all Docker containers using Portainer
   ```

3. **Expected behavior:**
   - Copilot should invoke the `portainer_list_containers` tool
   - Show container information in the chat

If tools work, the integration is successful!

## JetBrains IDEs (IntelliJ, PyCharm, WebStorm, etc.)

### Prerequisites

- **JetBrains IDE** (2023.3+ recommended)
- **GitHub Copilot plugin** installed and updated

### Configuration

1. **Open Settings:**
   - macOS: `Cmd+,`
   - Windows/Linux: `Ctrl+Alt+S`

2. **Navigate to:**
   - **Tools → GitHub Copilot → MCP Servers**

3. **Add MCP server:**
   - Click "Add Server"
   - **Name:** `axinova-tools`
   - **Command:** `/usr/local/bin/axinova-mcp-server`
   - **Environment Variables:** Add each variable

Example variables:
```
APP_PORTAINER__URL=https://portainer.axinova-internal.xyz
APP_PORTAINER__TOKEN=ptr_YOUR_TOKEN_HERE
APP_GRAFANA__URL=https://grafana.axinova-internal.xyz
APP_GRAFANA__TOKEN=glsa_YOUR_TOKEN_HERE
APP_TLS__SKIP_VERIFY=true
```

4. **Restart IDE**

5. **Test in Copilot Chat:**
   - Open Copilot Chat panel
   - Type: "Show all Docker containers"

**Note:** MCP support in JetBrains may vary by plugin version. Check the latest Copilot plugin documentation.

## GitHub Copilot CLI

### Prerequisites

- **GitHub CLI** (`gh`) installed
- **Copilot CLI extension** installed:
  ```bash
  gh extension install github/gh-copilot
  ```

### Configuration

The Copilot CLI reads MCP server configuration from environment variables.

**Option 1: Environment Variable**

```bash
export COPILOT_MCP_SERVERS='{
  "axinova-tools": {
    "command": "/usr/local/bin/axinova-mcp-server",
    "env": {
      "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
      "APP_PORTAINER__TOKEN": "ptr_YOUR_TOKEN_HERE",
      "APP_TLS__SKIP_VERIFY": "true"
    }
  }
}'
```

Add to your shell profile (`~/.bashrc`, `~/.zshrc`) to persist.

**Option 2: Config File**

Create `~/.config/github-copilot/mcp.json`:

```json
{
  "servers": {
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

**Note:** Check `gh copilot --help` for the exact config file location, as it may vary.

### Usage

```bash
# Suggest a command using MCP tools
gh copilot suggest "list all docker containers"

# Ask a question
gh copilot explain "what containers are running?"
```

Copilot CLI will use the MCP server to answer questions about your infrastructure.

## Enterprise Policy Setup

**For GitHub Copilot Business/Enterprise users:**

Organization administrators must enable MCP support before users can add MCP servers.

### Admin Steps

1. **Go to GitHub.com:**
   - Navigate to your **Organization**
   - Go to **Settings → Copilot**

2. **Enable MCP Policy:**
   - Find **"MCP servers in Copilot"** setting
   - Enable the policy

3. **Optional - Approve Specific Servers:**
   - Add `axinova-mcp-server` to allowlist
   - Set binary path restrictions

4. **Notify team members:**
   - Users can now configure MCP servers in their IDEs

### User Steps (After Policy Enabled)

Follow the VS Code or JetBrains setup instructions above.

**Important:** If you don't see MCP settings in your IDE, check with your GitHub organization admin to ensure the policy is enabled.

## Configuration Examples

### Minimal Configuration (Portainer Only)

**VS Code:**
```json
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

### Full Configuration (All Services)

**VS Code:**
```json
{
  "github.copilot.mcp.servers": {
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

### Docker-based Installation

**VS Code:**
```json
{
  "github.copilot.mcp.servers": {
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

### Team Project Configuration (Template)

For team projects, create a template without tokens:

**.vscode/settings.json.template:**
```json
{
  "github.copilot.mcp.servers": {
    "axinova-tools": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
        "APP_PORTAINER__TOKEN": "REPLACE_WITH_YOUR_TOKEN",
        "APP_GRAFANA__URL": "https://grafana.axinova-internal.xyz",
        "APP_GRAFANA__TOKEN": "REPLACE_WITH_YOUR_TOKEN"
      }
    }
  }
}
```

**README.md:**
```markdown
## Setup

1. Copy `.vscode/settings.json.template` to `.vscode/settings.json`
2. Replace `REPLACE_WITH_YOUR_TOKEN` with your actual tokens
3. Don't commit `.vscode/settings.json` (already in .gitignore)
```

**.gitignore:**
```
.vscode/settings.json
```

## Usage Examples

Once configured, you can interact with infrastructure using natural language in Copilot Chat:

### Container Management

**In Copilot Chat (Cmd+I / Ctrl+I):**
```
Show me all running containers

Get logs for the grafana container

What containers are consuming the most memory?

Restart the postgres container
```

### Monitoring

```
List all Grafana dashboards

What's on the "Container Resources" dashboard?

Show me active Prometheus alerts

What's the CPU usage trend for the last hour?
```

### Code Comments

You can also use MCP tools via code comments:

```python
# @copilot List all containers and filter by status=running


# @copilot Check if there are any Prometheus alerts firing
```

Copilot will invoke the tools and provide results as code comments.

### Inline Chat

In the editor, select code and press `Cmd+I` (macOS) or `Ctrl+I` (Windows):

```
Refactor this deployment script to check container health using Portainer before proceeding
```

Copilot will use the MCP tools to understand current container state.

### Chat Panel

In the Chat panel (sidebar), you can have multi-turn conversations:

```
You: What containers are running?
Copilot: [Lists containers using portainer_list_containers]

You: Show me logs for the first one
Copilot: [Gets logs using portainer_get_container_logs]

You: Create a task to investigate that error
Copilot: [Creates task using vikunja_create_task]
```

## Troubleshooting

### MCP Settings Not Available

**Problem:** Can't find "MCP servers" setting in VS Code

**Solution:**
1. **Update VS Code:**
   - Requires VS Code 1.102 or later
   - Help → Check for Updates

2. **Update GitHub Copilot extension:**
   - Extensions → GitHub Copilot → Update

3. **Check for preview features:**
   - Some MCP features may be in preview/experimental settings

4. **Enterprise policy:**
   - For Business/Enterprise, check with admin that MCP is enabled

### Tools Not Available in Copilot Chat

**Problem:** Copilot doesn't recognize MCP tools

**Solution:**
1. **Restart VS Code:**
   - Completely quit and reopen

2. **Check Developer Tools:**
   - Help → Toggle Developer Tools
   - Console tab → look for MCP-related errors

3. **Verify settings.json:**
   ```bash
   # macOS
   cat ~/Library/Application\ Support/Code/User/settings.json | jq .
   ```

   Check for JSON syntax errors.

4. **Test server manually:**
   ```bash
   echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | /usr/local/bin/axinova-mcp-server
   ```

   Should return JSON with server capabilities.

### Connection Errors to Services

**Problem:** Tools fail with HTTP errors

**Solution:**
1. **Verify service URLs:**
   ```bash
   curl -I https://portainer.axinova-internal.xyz
   curl -I https://grafana.axinova-internal.xyz
   ```

2. **Check tokens:**
   ```bash
   # Portainer
   curl -H "X-API-Key: ptr_YOUR_TOKEN" \
     https://portainer.axinova-internal.xyz/api/endpoints

   # Grafana
   curl -H "Authorization: Bearer glsa_YOUR_TOKEN" \
     https://grafana.axinova-internal.xyz/api/dashboards
   ```

3. **TLS verification:**
   - For self-signed certs: `"APP_TLS__SKIP_VERIFY": "true"`

### Permission Denied

**Problem:** Server fails to execute

**Solution:**
```bash
# Make binary executable
chmod +x /usr/local/bin/axinova-mcp-server

# Check ownership and permissions
ls -la /usr/local/bin/axinova-mcp-server
# Should show: -rwxr-xr-x
```

### Enterprise Policy Not Enabled

**Problem:** Can't add MCP servers (Enterprise users)

**Solution:**
1. **Contact your GitHub org admin**
2. **Ask them to:**
   - Go to Organization Settings → Copilot
   - Enable "MCP servers in Copilot"
3. **Restart VS Code after policy is enabled**

### Docker Container Issues

**Problem:** Docker-based MCP server fails to start

**Solution:**
1. **Test Docker command manually:**
   ```bash
   docker run -i --rm \
     -e APP_PORTAINER__TOKEN=ptr_xxx \
     ghcr.io/axinova-ai/axinova-mcp-server-go:latest
   ```

2. **Check Docker is running:**
   ```bash
   docker ps
   ```

3. **Pull latest image:**
   ```bash
   docker pull ghcr.io/axinova-ai/axinova-mcp-server-go:latest
   ```

## Security Best Practices

### 1. Protect Your Settings File

**VS Code User Settings:**
```bash
# Restrict permissions
chmod 600 ~/Library/Application\ Support/Code/User/settings.json
```

**Workspace Settings:**
- Add `.vscode/settings.json` to `.gitignore`
- Never commit tokens to version control

### 2. Use Environment Variables

Instead of hardcoding tokens, reference environment variables:

**.env (not committed):**
```bash
PORTAINER_TOKEN=ptr_xxx
GRAFANA_TOKEN=glsa_xxx
```

**.vscode/settings.json:**
```json
{
  "github.copilot.mcp.servers": {
    "axinova-tools": {
      "command": "/usr/local/bin/axinova-mcp-server",
      "env": {
        "APP_PORTAINER__TOKEN": "${env:PORTAINER_TOKEN}",
        "APP_GRAFANA__TOKEN": "${env:GRAFANA_TOKEN}"
      }
    }
  }
}
```

**Note:** Check if VS Code supports `${env:VAR}` syntax in Copilot MCP config.

### 3. Separate Dev and Prod

Configure multiple servers for different environments:

```json
{
  "github.copilot.mcp.servers": {
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

Specify which environment in Copilot Chat:
```
Using axinova-prod, list containers
```

### 4. Use Read-Only Tokens

Where possible, use tokens with minimal permissions:
- **Grafana:** Viewer role service accounts
- **Portainer:** Read-only API tokens

### 5. Rotate Tokens Regularly

Set reminders to rotate API tokens every 90 days.

## Known Limitations

### MCP Support is Preview

GitHub Copilot's MCP support is relatively new (added in 2026). Some limitations:

- May require latest VS Code / extension versions
- Enterprise policy required for Business/Enterprise users
- Tool execution may be slower than native Copilot features
- Not all MCP features may be supported

### Platform Availability

- **VS Code:** Full support (1.102+)
- **JetBrains:** Plugin-dependent, check latest docs
- **Copilot CLI:** Limited support, check `gh copilot` docs
- **GitHub.com:** No direct MCP support yet

### Network Access

MCP servers run locally, so they need network access to your services:
- Ensure firewall allows outbound HTTPS
- VPN may be required for internal services

## Updating the MCP Server

### Binary Installation

```bash
# Download new version
curl -L https://github.com/axinova-ai/axinova-mcp-server-go/releases/latest/download/axinova-mcp-server-macos -o /tmp/axinova-mcp-server

# Replace existing
sudo mv /tmp/axinova-mcp-server /usr/local/bin/axinova-mcp-server
sudo chmod +x /usr/local/bin/axinova-mcp-server

# Restart VS Code
```

### Docker Installation

```bash
# Pull latest image
docker pull ghcr.io/axinova-ai/axinova-mcp-server-go:latest

# Restart VS Code
```

Docker will use the new image on next run.

## Uninstalling

### Remove from VS Code

**User Settings:**
```bash
# Edit settings.json
code ~/Library/Application\ Support/Code/User/settings.json

# Remove the "github.copilot.mcp.servers" section
```

**Workspace Settings:**
```bash
# Delete or edit .vscode/settings.json
rm .vscode/settings.json
```

Restart VS Code.

### Remove Binary

```bash
sudo rm /usr/local/bin/axinova-mcp-server
```

### Remove Docker Image

```bash
docker rmi ghcr.io/axinova-ai/axinova-mcp-server-go:latest
```

## Next Steps

- **Explore tools:** See [NATIVE-MCP-INTEGRATION.md](../NATIVE-MCP-INTEGRATION.md) for full tool list
- **Try other clients:** [Claude Desktop](claude-desktop.md) or [Claude Code](claude-code.md)
- **Learn more:** [Example configurations](../examples/)

## Support

- **GitHub Issues:** Report bugs or request features
- **Documentation:** Full docs in this repository
- **GitHub Copilot Docs:** https://docs.github.com/copilot
- **MCP Spec:** https://spec.modelcontextprotocol.io/
