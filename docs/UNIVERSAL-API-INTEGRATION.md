# Universal API Integration Guide

**Version:** 1.0.0
**Last Updated:** 2026-01-25

---

## Overview

The Axinova MCP Server provides a **universal HTTP JSON-RPC API** that can be integrated with any LLM platform, not just Claude. This guide shows you how to integrate the 38 available tools with ChatGPT, Gemini, LangChain, LlamaIndex, and other AI platforms.

---

## Quick Start

### 1. Get API Access

**API Endpoint:**
```
https://mcp.axinova-ai.com
```

**Authentication:**
```
Bearer Token (contact admin for access)
```

**Available Tools:** 38 tools across 5 services
- **Portainer** (8 tools) - Docker container management
- **Grafana** (9 tools) - Monitoring dashboards
- **Prometheus** (7 tools) - Metrics and alerting
- **SilverBullet** (6 tools) - Wiki and knowledge base
- **Vikunja** (8 tools) - Task and project management

### 2. Test Connection

```bash
# List all available tools
curl https://mcp.axinova-ai.com/api/mcp/v1/tools \
  -H "Authorization: Bearer YOUR_TOKEN"

# Returns:
# {
#   "count": 38,
#   "tools": [...]
# }
```

### 3. Call a Tool

```bash
curl -X POST https://mcp.axinova-ai.com/api/mcp/v1/call \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "grafana_list_dashboards",
      "arguments": {}
    }
  }'

# Returns dashboard list
```

---

## API Reference

### Authentication

All requests require a Bearer token in the Authorization header:

```
Authorization: Bearer YOUR_TOKEN_HERE
```

### Endpoints

#### 1. List Tools

**Endpoint:** `GET /api/mcp/v1/tools`

**Response:**
```json
{
  "count": 38,
  "tools": [
    {
      "name": "portainer_list_containers",
      "description": "List all Docker containers in a Portainer environment",
      "inputSchema": {
        "type": "object",
        "properties": {
          "endpoint_id": {
            "type": "number",
            "description": "Portainer endpoint ID (default: 1 for local)"
          }
        }
      }
    },
    ...
  ]
}
```

#### 2. Call Tool

**Endpoint:** `POST /api/mcp/v1/call`

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "tool_name",
    "arguments": {
      "arg1": "value1",
      "arg2": "value2"
    }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "content": [
      {
        "type": "text",
        "text": "Tool execution result..."
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
    "code": -32000,
    "message": "Tool not found",
    "data": "tool_name"
  }
}
```

### Rate Limiting

- **Rate:** 60 requests per minute per IP
- **Burst:** 10 requests
- **Header:** `X-RateLimit-Remaining`

### Best Practices

1. **Cache tool list** - Fetch once and reuse
2. **Handle errors** - Always check for `error` field in response
3. **Use connection pooling** - Reuse HTTP connections
4. **Set timeouts** - Default timeout: 30 seconds
5. **Respect rate limits** - Check `X-RateLimit-Remaining` header

---

## Integration Guides

### Available Platform Guides

- **[ChatGPT Integration](integrations/chatgpt.md)** - OpenAI Function Calling
- **[Gemini Integration](integrations/gemini.md)** - Google AI Function Calling
- **[LangChain Integration](integrations/langchain.md)** - Custom Tools
- **[LlamaIndex Integration](integrations/llamaindex.md)** - Tool Integration

### Quick Integration Example (Python)

```python
import requests
import json

MCP_API_URL = "https://mcp.axinova-ai.com"
MCP_TOKEN = "your-token-here"

def call_mcp_tool(tool_name, arguments=None):
    """Call an MCP tool and return the result."""
    response = requests.post(
        f"{MCP_API_URL}/api/mcp/v1/call",
        headers={
            "Authorization": f"Bearer {MCP_TOKEN}",
            "Content-Type": "application/json"
        },
        json={
            "jsonrpc": "2.0",
            "id": 1,
            "method": "tools/call",
            "params": {
                "name": tool_name,
                "arguments": arguments or {}
            }
        },
        timeout=30
    )

    response.raise_for_status()
    data = response.json()

    if "error" in data:
        raise Exception(f"MCP Error: {data['error']['message']}")

    return data["result"]["content"][0]["text"]

# Example usage
dashboards = call_mcp_tool("grafana_list_dashboards")
print(dashboards)

containers = call_mcp_tool("portainer_list_containers", {"endpoint_id": 1})
print(containers)
```

---

## Tool Discovery

### Fetch and Convert Tool Schemas

```python
import requests

def get_mcp_tools():
    """Fetch all available MCP tools."""
    response = requests.get(
        "https://mcp.axinova-ai.com/api/mcp/v1/tools",
        headers={"Authorization": f"Bearer {YOUR_TOKEN}"}
    )
    return response.json()["tools"]

def convert_to_openai_function(tool):
    """Convert MCP tool schema to OpenAI function format."""
    return {
        "type": "function",
        "function": {
            "name": tool["name"],
            "description": tool["description"],
            "parameters": tool["inputSchema"]
        }
    }

# Get all tools and convert
mcp_tools = get_mcp_tools()
openai_functions = [convert_to_openai_function(t) for t in mcp_tools]
```

---

## Available Tools

### Portainer (8 tools)

- `portainer_list_containers` - List all Docker containers
- `portainer_start_container` - Start a container
- `portainer_stop_container` - Stop a container
- `portainer_restart_container` - Restart a container
- `portainer_get_container_logs` - Get container logs
- `portainer_list_stacks` - List Docker Compose stacks
- `portainer_get_stack` - Get stack details
- `portainer_inspect_container` - Inspect container details

### Grafana (9 tools)

- `grafana_list_dashboards` - List all dashboards
- `grafana_get_dashboard` - Get dashboard by UID
- `grafana_create_dashboard` - Create new dashboard
- `grafana_delete_dashboard` - Delete dashboard
- `grafana_list_datasources` - List all datasources
- `grafana_create_datasource` - Create datasource
- `grafana_query_datasource` - Query datasource
- `grafana_list_alert_rules` - List alert rules
- `grafana_get_health` - Check Grafana health

### Prometheus (7 tools)

- `prometheus_query` - Execute instant query
- `prometheus_query_range` - Execute range query
- `prometheus_list_label_names` - List all label names
- `prometheus_list_label_values` - List values for label
- `prometheus_find_series` - Find time series
- `prometheus_list_targets` - List scrape targets
- `prometheus_get_metadata` - Get metric metadata

### SilverBullet (6 tools)

- `silverbullet_list_pages` - List all pages
- `silverbullet_get_page` - Get page content
- `silverbullet_create_page` - Create new page
- `silverbullet_update_page` - Update page
- `silverbullet_delete_page` - Delete page
- `silverbullet_search_pages` - Search pages

### Vikunja (8 tools)

- `vikunja_list_projects` - List all projects
- `vikunja_get_project` - Get project details
- `vikunja_create_project` - Create project
- `vikunja_list_tasks` - List tasks in project
- `vikunja_get_task` - Get task details
- `vikunja_create_task` - Create task
- `vikunja_update_task` - Update task
- `vikunja_delete_task` - Delete task

---

## Error Handling

### Common Error Codes

| Code | Meaning | Solution |
|------|---------|----------|
| 401 | Unauthorized | Check Bearer token |
| 404 | Not Found | Tool name incorrect |
| 429 | Rate Limited | Wait before retry |
| -32000 | Tool Error | Check tool arguments |
| -32601 | Method Not Found | Invalid method name |
| -32602 | Invalid Params | Check argument types |
| -32700 | Parse Error | Invalid JSON |

### Error Handling Example

```python
def safe_call_tool(tool_name, arguments=None):
    """Call MCP tool with error handling."""
    try:
        response = requests.post(
            f"{MCP_API_URL}/api/mcp/v1/call",
            headers={
                "Authorization": f"Bearer {MCP_TOKEN}",
                "Content-Type": "application/json"
            },
            json={
                "jsonrpc": "2.0",
                "id": 1,
                "method": "tools/call",
                "params": {
                    "name": tool_name,
                    "arguments": arguments or {}
                }
            },
            timeout=30
        )

        # Check HTTP status
        if response.status_code == 401:
            raise Exception("Invalid API token")
        elif response.status_code == 429:
            raise Exception("Rate limit exceeded, retry later")

        response.raise_for_status()
        data = response.json()

        # Check JSON-RPC error
        if "error" in data:
            error = data["error"]
            raise Exception(f"Tool error ({error['code']}): {error['message']}")

        return data["result"]["content"][0]["text"]

    except requests.Timeout:
        raise Exception("Request timeout (30s)")
    except requests.ConnectionError:
        raise Exception("Connection failed")
    except Exception as e:
        raise Exception(f"MCP API Error: {str(e)}")
```

---

## Example Use Cases

### Use Case 1: Monitor Container Health

```python
# List all containers
containers = call_mcp_tool("portainer_list_containers", {"endpoint_id": 1})
containers_list = json.loads(containers)

# Check for unhealthy containers
for container in containers_list:
    if container["State"] != "running":
        print(f"Warning: {container['Names']} is {container['State']}")

        # Get logs
        logs = call_mcp_tool("portainer_get_container_logs", {
            "container_id": container["Id"],
            "tail": 50
        })
        print(f"Logs:\n{logs}")
```

### Use Case 2: Dashboard Analytics

```python
# Get all dashboards
dashboards = call_mcp_tool("grafana_list_dashboards")
dashboard_list = json.loads(dashboards)

print(f"Total dashboards: {len(dashboard_list)}")

# Get details for each dashboard
for dashboard in dashboard_list:
    details = call_mcp_tool("grafana_get_dashboard", {"uid": dashboard["uid"]})
    print(f"Dashboard: {dashboard['title']}")
    print(f"Tags: {dashboard['tags']}")
```

### Use Case 3: Metric Queries

```python
# Query CPU usage
cpu_usage = call_mcp_tool("prometheus_query", {
    "query": "rate(container_cpu_usage_seconds_total[5m])"
})
print(f"CPU Usage:\n{cpu_usage}")

# Query memory usage range
memory_range = call_mcp_tool("prometheus_query_range", {
    "query": "container_memory_usage_bytes",
    "start": "1h",
    "step": "1m"
})
print(f"Memory Range:\n{memory_range}")
```

### Use Case 4: Task Management

```python
# List all projects
projects = call_mcp_tool("vikunja_list_projects")
project_list = json.loads(projects)

if project_list:
    project_id = project_list[0]["id"]

    # Create a task
    task = call_mcp_tool("vikunja_create_task", {
        "project_id": project_id,
        "title": "Review monitoring alerts",
        "description": "Check Grafana dashboards",
        "priority": 3
    })
    print(f"Created task: {task}")
```

---

## Support

### Documentation
- [ChatGPT Integration](integrations/chatgpt.md)
- [Gemini Integration](integrations/gemini.md)
- [LangChain Integration](integrations/langchain.md)
- [LlamaIndex Integration](integrations/llamaindex.md)

### API Status
- **Uptime:** Check https://mcp.axinova-ai.com:9001/health
- **Metrics:** Check https://mcp.axinova-ai.com:9001/metrics

### Contact
- **Issues:** https://github.com/axinova-ai/axinova-mcp-server-go/issues
- **Email:** [Contact your administrator]

---

## License

API access is subject to terms of service. Contact administrator for commercial use.
