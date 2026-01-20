# Testing the MCP Server

This guide provides detailed instructions for testing the Axinova MCP Server with various MCP clients.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Testing with Claude Desktop](#testing-with-claude-desktop)
- [Testing with Python MCP SDK](#testing-with-python-mcp-sdk)
- [Testing with Node.js MCP SDK](#testing-with-nodejs-mcp-sdk)
- [Manual Testing with stdio](#manual-testing-with-stdio)
- [Testing Individual Tools](#testing-individual-tools)
- [Troubleshooting](#troubleshooting)

---

## Prerequisites

1. **MCP Server Built and Running**
   ```bash
   cd axinova-mcp-server-go
   make build
   ```

2. **Environment Variables Configured**
   - Copy `.env.example` to `.env`
   - Fill in API tokens for each service
   - See [scripts/get-tokens.md](scripts/get-tokens.md) for token generation

3. **Network Access**
   - Ensure connectivity to internal services
   - VPN connection if accessing from outside network

---

## Testing with Claude Desktop

### 1. Install Claude Desktop

Download from: https://claude.ai/download

### 2. Configure MCP Server

**macOS/Linux:**
Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "axinova": {
      "command": "/Users/weixia/axinova/axinova-mcp-server-go/bin/axinova-mcp-server",
      "env": {
        "ENV": "prod",
        "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
        "APP_PORTAINER__TOKEN": "your-portainer-token",
        "APP_GRAFANA__URL": "https://grafana.axinova-internal.xyz",
        "APP_GRAFANA__TOKEN": "your-grafana-token",
        "APP_PROMETHEUS__URL": "https://prometheus.axinova-internal.xyz",
        "APP_SILVERBULLET__URL": "https://silverbullet.axinova-internal.xyz",
        "APP_SILVERBULLET__TOKEN": "your-silverbullet-token",
        "APP_VIKUNJA__URL": "https://vikunja.axinova-internal.xyz",
        "APP_VIKUNJA__TOKEN": "your-vikunja-token"
      }
    }
  }
}
```

**Windows:**
Edit `%APPDATA%\Claude\claude_desktop_config.json`

### 3. Restart Claude Desktop

After saving the config, restart Claude Desktop completely.

### 4. Verify Connection

In Claude Desktop, look for the MCP indicator (hammer icon 🔨) in the bottom right. Click it to see available tools.

### 5. Test Commands

Try these example prompts in Claude:

```
List all Docker containers in Portainer
```

```
Show me Grafana dashboards
```

```
Query Prometheus for CPU usage: up
```

```
List all my Vikunja projects
```

```
Create a new SilverBullet page called "Test Note" with content "This is a test"
```

### 6. Expected Behavior

- Claude should show which tool it's using (e.g., "Using portainer_list_containers")
- Results should be displayed in the conversation
- Error messages should be descriptive

---

## Testing with Python MCP SDK

### 1. Install Python MCP SDK

```bash
pip install mcp
```

### 2. Create Test Script

Save as `test_mcp_client.py`:

```python
#!/usr/bin/env python3
import asyncio
import json
from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client

async def main():
    server_params = StdioServerParameters(
        command="/Users/weixia/axinova/axinova-mcp-server-go/bin/axinova-mcp-server",
        env={
            "ENV": "prod",
            "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
            "APP_PORTAINER__TOKEN": "your-token",
            # ... add other env vars
        }
    )

    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            # Initialize
            await session.initialize()

            # List available tools
            tools = await session.list_tools()
            print(f"Available tools: {len(tools.tools)}")
            for tool in tools.tools[:5]:
                print(f"  - {tool.name}: {tool.description}")

            # Call a tool
            result = await session.call_tool("prometheus_query", {"query": "up"})
            print(f"\nPrometheus query result:")
            print(json.dumps(result.content, indent=2))

if __name__ == "__main__":
    asyncio.run(main())
```

### 3. Run Test

```bash
python3 test_mcp_client.py
```

---

## Testing with Node.js MCP SDK

### 1. Install Node.js MCP SDK

```bash
npm install @modelcontextprotocol/sdk
```

### 2. Create Test Script

Save as `test_mcp_client.js`:

```javascript
#!/usr/bin/env node

import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";

async function main() {
  const transport = new StdioClientTransport({
    command: "/Users/weixia/axinova/axinova-mcp-server-go/bin/axinova-mcp-server",
    env: {
      ENV: "prod",
      APP_PORTAINER__URL: "https://portainer.axinova-internal.xyz",
      APP_PORTAINER__TOKEN: "your-token",
      // ... add other env vars
    }
  });

  const client = new Client({
    name: "test-client",
    version: "1.0.0"
  }, {
    capabilities: {}
  });

  await client.connect(transport);

  // List tools
  const tools = await client.listTools();
  console.log(`Available tools: ${tools.tools.length}`);
  tools.tools.slice(0, 5).forEach(tool => {
    console.log(`  - ${tool.name}: ${tool.description}`);
  });

  // Call a tool
  const result = await client.callTool({
    name: "prometheus_query",
    arguments: { query: "up" }
  });
  console.log("\nPrometheus query result:");
  console.log(JSON.stringify(result, null, 2));

  await client.close();
}

main().catch(console.error);
```

### 3. Run Test

```bash
node test_mcp_client.js
```

---

## Manual Testing with stdio

### 1. Interactive Testing

Use the provided test script:

```bash
./test_mcp.sh
```

This will:
- Send initialize request
- List all available tools
- Show tool schemas

### 2. Custom JSON-RPC Commands

Manually send commands via stdio:

```bash
# Initialize
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}' | ./bin/axinova-mcp-server

# List tools
(echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}'; \
 echo '{"jsonrpc":"2.0","method":"initialized"}'; \
 echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}') | ./bin/axinova-mcp-server 2>/dev/null | grep -v '^\[MCP\]'

# Call a tool
(echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}'; \
 echo '{"jsonrpc":"2.0","method":"initialized"}'; \
 echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"prometheus_query","arguments":{"query":"up"}}}') | ./bin/axinova-mcp-server
```

---

## Testing Individual Tools

### Portainer Tools

```json
// List containers
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "portainer_list_containers",
    "arguments": {
      "endpoint_id": 1
    }
  }
}

// Get container logs
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "portainer_get_container_logs",
    "arguments": {
      "container_id": "your-container-id",
      "tail": 50
    }
  }
}
```

### Grafana Tools

```json
// List dashboards
{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "tools/call",
  "params": {
    "name": "grafana_list_dashboards",
    "arguments": {}
  }
}

// Query datasource
{
  "jsonrpc": "2.0",
  "id": 7,
  "method": "tools/call",
  "params": {
    "name": "grafana_query_datasource",
    "arguments": {
      "datasource_uid": "your-datasource-uid",
      "query": "up"
    }
  }
}
```

### Prometheus Tools

```json
// Instant query
{
  "jsonrpc": "2.0",
  "id": 8,
  "method": "tools/call",
  "params": {
    "name": "prometheus_query",
    "arguments": {
      "query": "rate(http_requests_total[5m])"
    }
  }
}

// Range query
{
  "jsonrpc": "2.0",
  "id": 9,
  "method": "tools/call",
  "params": {
    "name": "prometheus_query_range",
    "arguments": {
      "query": "up",
      "start": "1h",
      "step": "15s"
    }
  }
}

// List targets
{
  "jsonrpc": "2.0",
  "id": 10,
  "method": "tools/call",
  "params": {
    "name": "prometheus_list_targets",
    "arguments": {}
  }
}
```

### SilverBullet Tools

```json
// List pages
{
  "jsonrpc": "2.0",
  "id": 11,
  "method": "tools/call",
  "params": {
    "name": "silverbullet_list_pages",
    "arguments": {}
  }
}

// Create page
{
  "jsonrpc": "2.0",
  "id": 12,
  "method": "tools/call",
  "params": {
    "name": "silverbullet_create_page",
    "arguments": {
      "page_name": "Test Page",
      "content": "# Test\\n\\nThis is a test page created via MCP."
    }
  }
}
```

### Vikunja Tools

```json
// List projects
{
  "jsonrpc": "2.0",
  "id": 13,
  "method": "tools/call",
  "params": {
    "name": "vikunja_list_projects",
    "arguments": {}
  }
}

// Create task
{
  "jsonrpc": "2.0",
  "id": 14,
  "method": "tools/call",
  "params": {
    "name": "vikunja_create_task",
    "arguments": {
      "project_id": 1,
      "title": "Test MCP Server",
      "description": "Testing task creation via MCP",
      "priority": 3
    }
  }
}
```

---

## Troubleshooting

### Server Not Starting

**Issue:** Server fails to start

**Solutions:**
1. Check environment variables are set correctly
2. Verify binary has execute permissions: `chmod +x bin/axinova-mcp-server`
3. Check logs in stderr: `./bin/axinova-mcp-server 2>&1`

### TLS Certificate Errors

**Issue:** `x509: certificate signed by unknown authority`

**Solution:** Ensure `config/prod.yaml` has `tls.skip_verify: true`

### Authentication Errors

**Issue:** HTTP 401 Unauthorized

**Solutions:**
1. Verify API tokens are correct and not expired
2. Check token format (Bearer vs API Key)
3. Test token manually with curl:
   ```bash
   curl -H "Authorization: Bearer YOUR_TOKEN" https://service/api/endpoint
   ```

### Tool Not Found

**Issue:** "Tool not found" error

**Solution:**
1. List available tools: send `tools/list` request
2. Check tool name spelling
3. Verify service is enabled in config

### Connection Timeout

**Issue:** HTTP timeouts

**Solutions:**
1. Check network connectivity: `ping portainer.axinova-internal.xyz`
2. Increase timeout in `config/base.yaml`:
   ```yaml
   timeout:
     http: 60s
   ```
3. Verify services are running

### Claude Desktop Not Detecting Tools

**Issue:** No tools appear in Claude Desktop

**Solutions:**
1. Check `claude_desktop_config.json` syntax
2. Restart Claude Desktop completely (not just reload)
3. Check MCP server logs:
   ```bash
   tail -f ~/Library/Logs/Claude/mcp*.log
   ```
4. Verify server binary path is absolute

### Empty Results

**Issue:** Tool returns empty results

**Solutions:**
1. Verify the service has data (e.g., containers exist, dashboards exist)
2. Check API token permissions
3. Test API endpoint directly with curl

---

## Debugging Tips

### Enable Debug Logging

Set environment variable:
```bash
export APP_LOG__LEVEL=debug
./bin/axinova-mcp-server
```

### Capture Full stdio Exchange

```bash
./bin/axinova-mcp-server < commands.json > responses.json 2> errors.log
```

### Test with Minimal Config

Disable all services except one:
```bash
export APP_GRAFANA__ENABLED=false
export APP_PORTAINER__ENABLED=false
export APP_SILVERBULLET__ENABLED=false
export APP_VIKUNJA__ENABLED=false
# Only Prometheus enabled
./bin/axinova-mcp-server
```

---

## Example Testing Session

Complete working example:

```bash
#!/bin/bash

# 1. Set environment
export ENV=prod
export APP_PORTAINER__URL=https://portainer.axinova-internal.xyz
export APP_PORTAINER__TOKEN=your-token
export APP_PROMETHEUS__URL=https://prometheus.axinova-internal.xyz

# 2. Build server
make build

# 3. Test initialization
echo "Testing initialization..."
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}' | \
  ./bin/axinova-mcp-server 2>&1 | grep -q "protocolVersion" && echo "✓ Init OK" || echo "✗ Init failed"

# 4. Test tool listing
echo "Testing tool listing..."
(echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}'; \
 echo '{"jsonrpc":"2.0","method":"initialized"}'; \
 echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}') | \
  ./bin/axinova-mcp-server 2>&1 | grep -q "prometheus_query" && echo "✓ Tools OK" || echo "✗ Tools failed"

# 5. Test actual query
echo "Testing Prometheus query..."
(echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}'; \
 echo '{"jsonrpc":"2.0","method":"initialized"}'; \
 echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"prometheus_query","arguments":{"query":"up"}}}') | \
  ./bin/axinova-mcp-server 2>&1 | grep -q '"status":"success"' && echo "✓ Query OK" || echo "✗ Query failed"

echo "Testing complete!"
```

---

## Performance Testing

### Load Testing

Test with multiple concurrent requests:

```python
import asyncio
from mcp import ClientSession, StdioServerParameters
from mcp.client.stdio import stdio_client

async def test_concurrent():
    server_params = StdioServerParameters(
        command="./bin/axinova-mcp-server",
        env={"ENV": "prod"}
    )

    async with stdio_client(server_params) as (read, write):
        async with ClientSession(read, write) as session:
            await session.initialize()

            # Run 10 queries concurrently
            tasks = [
                session.call_tool("prometheus_query", {"query": "up"})
                for _ in range(10)
            ]

            results = await asyncio.gather(*tasks)
            print(f"Completed {len(results)} concurrent requests")

asyncio.run(test_concurrent())
```

---

## Next Steps

After successful testing:

1. **Deploy to Production** - See [DEPLOYMENT.md](DEPLOYMENT.md)
2. **Monitor Usage** - Check server logs and performance
3. **Add More Tools** - Extend with additional services
4. **Create Workflows** - Build complex automation with tool combinations
