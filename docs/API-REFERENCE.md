# MCP Server API Reference

## Overview

The Axinova MCP Server provides HTTP JSON-RPC access to productivity tools via a unified API. The server implements the Model Context Protocol (MCP) specification (2025-11-25) and exposes both native MCP tools and HTTP JSON-RPC endpoints.

**Base URL:** `https://mcp.axinova-ai.com`

## Authentication

All API requests require Bearer token authentication via the `Authorization` header:

```bash
Authorization: Bearer sk-mcp-prod-86d850ac73a8b9dd11e94b104ea4fd56966bee365ed5ffa3820ecd99f5f2640e
```

## Endpoints

### GET /health

Health check endpoint for monitoring.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2026-01-24T10:30:00Z"
}
```

### GET /ready

Readiness probe for container orchestration.

**Response:**
```json
{
  "status": "ready"
}
```

### GET /status

Server status and version information.

**Response:**
```json
{
  "status": "running",
  "version": "1.0.0",
  "uptime": "1h30m",
  "mcp_protocol_version": "2025-11-25"
}
```

### GET /metrics

Prometheus metrics endpoint (no authentication required).

**Exposed Metrics:**
- `mcp_server_uptime_seconds` - Server uptime gauge
- `mcp_tools_registered_total` - Number of registered tools
- `mcp_resources_registered_total` - Number of registered resources
- `mcp_rpc_requests_total{method, transport}` - Total RPC requests counter
- `mcp_rpc_request_duration_seconds{method, transport}` - Request duration histogram
- `mcp_rpc_errors_total{method, error_code, transport}` - Error counter
- `mcp_http_active_connections` - Active HTTP connections gauge

**Example:**
```bash
curl https://mcp.axinova-ai.com/metrics
```

### GET /api/mcp/v1/tools

List all available MCP tools with their schemas.

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "tools": [
    {
      "name": "portainer_list_containers",
      "description": "List all Docker containers in Portainer",
      "inputSchema": {
        "type": "object",
        "properties": {
          "endpoint_id": {
            "type": "number",
            "description": "Portainer endpoint ID"
          }
        },
        "required": ["endpoint_id"]
      }
    },
    {
      "name": "grafana_list_dashboards",
      "description": "List all Grafana dashboards",
      "inputSchema": {
        "type": "object",
        "properties": {}
      }
    }
  ],
  "count": 25
}
```

### POST /api/mcp/v1/call

Execute an MCP tool via JSON-RPC 2.0.

**Headers:**
```
Authorization: Bearer <token>
Content-Type: application/json
```

**Request Body:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "portainer_list_containers",
    "arguments": {
      "endpoint_id": 1
    }
  }
}
```

**Success Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "[\n  {\n    \"Id\": \"abc123\",\n    \"Names\": [\"/my-container\"],\n    \"State\": \"running\"\n  }\n]"
      }
    ]
  }
}
```

**Error Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32602,
    "message": "Invalid params: missing required field 'endpoint_id'"
  }
}
```

## Usage Examples

### cURL Examples

#### List Available Tools
```bash
curl https://mcp.axinova-ai.com/api/mcp/v1/tools \
  -H "Authorization: Bearer sk-mcp-prod-86d850ac73a8b9dd11e94b104ea4fd56966bee365ed5ffa3820ecd99f5f2640e"
```

#### Call a Tool
```bash
curl -X POST https://mcp.axinova-ai.com/api/mcp/v1/call \
  -H "Authorization: Bearer sk-mcp-prod-86d850ac73a8b9dd11e94b104ea4fd56966bee365ed5ffa3820ecd99f5f2640e" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "portainer_list_containers",
      "arguments": {
        "endpoint_id": 1
      }
    }
  }'
```

#### Health Check
```bash
curl https://mcp.axinova-ai.com/health
```

#### Fetch Metrics
```bash
curl https://mcp.axinova-ai.com/metrics | grep mcp_
```

### Python Client

```python
import requests
import json

class MCPClient:
    def __init__(self, api_url, api_token):
        self.api_url = api_url
        self.headers = {
            "Authorization": f"Bearer {api_token}",
            "Content-Type": "application/json"
        }

    def list_tools(self):
        """List all available MCP tools"""
        response = requests.get(
            f"{self.api_url}/api/mcp/v1/tools",
            headers=self.headers
        )
        response.raise_for_status()
        return response.json()

    def call_tool(self, tool_name, arguments):
        """Execute an MCP tool"""
        payload = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "tools/call",
            "params": {
                "name": tool_name,
                "arguments": arguments
            }
        }
        response = requests.post(
            f"{self.api_url}/api/mcp/v1/call",
            headers=self.headers,
            json=payload
        )
        response.raise_for_status()
        data = response.json()

        if "error" in data:
            raise Exception(f"MCP Error: {data['error']}")

        return data["result"]

# Usage
mcp = MCPClient(
    "https://mcp.axinova-ai.com",
    "sk-mcp-prod-86d850ac73a8b9dd11e94b104ea4fd56966bee365ed5ffa3820ecd99f5f2640e"
)

# List tools
tools = mcp.list_tools()
print(f"Available tools: {tools['count']}")

# Call a tool
result = mcp.call_tool("portainer_list_containers", {"endpoint_id": 1})
print(result)
```

### JavaScript/Node.js Client

```javascript
const axios = require('axios');

class MCPClient {
  constructor(apiUrl, apiToken) {
    this.apiUrl = apiUrl;
    this.headers = {
      'Authorization': `Bearer ${apiToken}`,
      'Content-Type': 'application/json'
    };
  }

  async listTools() {
    const response = await axios.get(
      `${this.apiUrl}/api/mcp/v1/tools`,
      { headers: this.headers }
    );
    return response.data;
  }

  async callTool(toolName, args) {
    const payload = {
      jsonrpc: '2.0',
      id: 1,
      method: 'tools/call',
      params: {
        name: toolName,
        arguments: args
      }
    };

    const response = await axios.post(
      `${this.apiUrl}/api/mcp/v1/call`,
      payload,
      { headers: this.headers }
    );

    if (response.data.error) {
      throw new Error(`MCP Error: ${JSON.stringify(response.data.error)}`);
    }

    return response.data.result;
  }
}

// Usage
const mcp = new MCPClient(
  'https://mcp.axinova-ai.com',
  'sk-mcp-prod-86d850ac73a8b9dd11e94b104ea4fd56966bee365ed5ffa3820ecd99f5f2640e'
);

// List tools
mcp.listTools().then(tools => {
  console.log(`Available tools: ${tools.count}`);
});

// Call a tool
mcp.callTool('portainer_list_containers', { endpoint_id: 1 })
  .then(result => console.log(result))
  .catch(err => console.error(err));
```

## Available Tools

See [TOOL-CATALOG.md](./TOOL-CATALOG.md) for a complete list of all available tools organized by service.

**Service Categories:**
- **Portainer** - Docker container management
- **Grafana** - Dashboard and monitoring management
- **Prometheus** - Metrics querying and analysis
- **SilverBullet** - Wiki/knowledge base management
- **Vikunja** - Task and project management

## Error Codes

The API uses standard JSON-RPC 2.0 error codes:

| Code | Message | Description |
|------|---------|-------------|
| -32700 | Parse error | Invalid JSON received |
| -32600 | Invalid request | JSON-RPC request is not valid |
| -32601 | Method not found | Method does not exist |
| -32602 | Invalid params | Invalid method parameters |
| -32603 | Internal error | Internal JSON-RPC error |
| -32000 | Server error | Generic server error (check message) |
| -32001 | Unauthorized | Missing or invalid API token |
| -32002 | Service unavailable | Backing service (Portainer, Grafana, etc.) is unavailable |

## Rate Limits

- **Default:** 1000 requests per minute per API token
- **Burst:** Up to 100 requests in a 10-second window
- Rate limit headers are included in responses:
  - `X-RateLimit-Limit`: Maximum requests per minute
  - `X-RateLimit-Remaining`: Requests remaining in current window
  - `X-RateLimit-Reset`: Unix timestamp when the limit resets

## Security Considerations

1. **Token Storage**: Store API tokens securely (environment variables, secrets manager)
2. **HTTPS Only**: All requests must use HTTPS
3. **Token Rotation**: Rotate API tokens every 90 days (recommended)
4. **Network Access**: API is exposed publicly but backing services are on private network
5. **Audit Logging**: All API calls are logged with timestamps and request details

## Support

For issues or questions:
- GitHub Issues: https://github.com/axinova-ai/axinova-mcp-server-go/issues
- Documentation: https://github.com/axinova-ai/axinova-mcp-server-go/tree/main/docs
