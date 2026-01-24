# MCP Server Tool Catalog

This document lists all available tools provided by the MCP server, organized by service category.

**Last Updated:** 2026-01-24

**Server URL:** `https://mcp.axinova-ai.com`

## Table of Contents

- [Portainer Tools](#portainer-tools) - Docker container management
- [Grafana Tools](#grafana-tools) - Dashboard and monitoring management
- [Prometheus Tools](#prometheus-tools) - Metrics querying and analysis
- [SilverBullet Tools](#silverbullet-tools) - Wiki and knowledge base management
- [Vikunja Tools](#vikunja-tools) - Task and project management

---

## Portainer Tools

Portainer provides Docker container and stack management capabilities.

**Service:** Portainer (ax-tools)
**URL:** `https://portainer.axinova-internal.xyz`

### portainer_list_containers

List all Docker containers in a Portainer endpoint.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "endpoint_id": {
      "type": "number",
      "description": "Portainer endpoint ID (usually 1 for local)"
    }
  },
  "required": ["endpoint_id"]
}
```

**Example:**
```bash
curl -X POST https://mcp.axinova-ai.com/api/mcp/v1/call \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "portainer_list_containers",
      "arguments": {"endpoint_id": 1}
    }
  }'
```

**Output:** JSON array of container objects with ID, Names, Image, State, Status, etc.

### portainer_list_stacks

List all Docker Compose stacks in a Portainer endpoint.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "endpoint_id": {
      "type": "number",
      "description": "Portainer endpoint ID"
    }
  },
  "required": ["endpoint_id"]
}
```

**Example:**
```bash
curl -X POST https://mcp.axinova-ai.com/api/mcp/v1/call \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "portainer_list_stacks",
      "arguments": {"endpoint_id": 1}
    }
  }'
```

**Output:** JSON array of stack objects with ID, Name, Type, Status, etc.

### portainer_get_container_logs

Get logs from a specific Docker container.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "endpoint_id": {
      "type": "number",
      "description": "Portainer endpoint ID"
    },
    "container_id": {
      "type": "string",
      "description": "Docker container ID or name"
    },
    "tail": {
      "type": "number",
      "description": "Number of lines to return from the end (default: 100)"
    }
  },
  "required": ["endpoint_id", "container_id"]
}
```

**Example:**
```bash
curl -X POST https://mcp.axinova-ai.com/api/mcp/v1/call \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "portainer_get_container_logs",
      "arguments": {
        "endpoint_id": 1,
        "container_id": "my-container",
        "tail": 50
      }
    }
  }'
```

**Output:** Container log output as text.

### portainer_inspect_container

Get detailed information about a specific container.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "endpoint_id": {
      "type": "number",
      "description": "Portainer endpoint ID"
    },
    "container_id": {
      "type": "string",
      "description": "Docker container ID or name"
    }
  },
  "required": ["endpoint_id", "container_id"]
}
```

**Output:** Detailed container configuration and state information.

### portainer_start_container

Start a stopped container.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "endpoint_id": {
      "type": "number",
      "description": "Portainer endpoint ID"
    },
    "container_id": {
      "type": "string",
      "description": "Docker container ID or name"
    }
  },
  "required": ["endpoint_id", "container_id"]
}
```

**Output:** Success confirmation or error message.

### portainer_stop_container

Stop a running container.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "endpoint_id": {
      "type": "number",
      "description": "Portainer endpoint ID"
    },
    "container_id": {
      "type": "string",
      "description": "Docker container ID or name"
    }
  },
  "required": ["endpoint_id", "container_id"]
}
```

**Output:** Success confirmation or error message.

### portainer_restart_container

Restart a container.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "endpoint_id": {
      "type": "number",
      "description": "Portainer endpoint ID"
    },
    "container_id": {
      "type": "string",
      "description": "Docker container ID or name"
    }
  },
  "required": ["endpoint_id", "container_id"]
}
```

**Output:** Success confirmation or error message.

---

## Grafana Tools

Grafana provides dashboard management and monitoring visualization.

**Service:** Grafana (ax-tools)
**URL:** `https://grafana.axinova-internal.xyz`

### grafana_list_dashboards

List all Grafana dashboards.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {}
}
```

**Example:**
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
```

**Output:** JSON array of dashboard objects with UID, title, URL, tags, etc.

### grafana_get_dashboard

Get a specific dashboard by UID.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "uid": {
      "type": "string",
      "description": "Dashboard UID"
    }
  },
  "required": ["uid"]
}
```

**Output:** Complete dashboard JSON including panels, queries, and configuration.

### grafana_create_dashboard

Create a new Grafana dashboard.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "title": {
      "type": "string",
      "description": "Dashboard title"
    },
    "dashboard": {
      "type": "object",
      "description": "Dashboard JSON configuration"
    },
    "folder_id": {
      "type": "number",
      "description": "Folder ID to save dashboard in (optional)"
    }
  },
  "required": ["title", "dashboard"]
}
```

**Output:** Created dashboard UID and URL.

### grafana_update_dashboard

Update an existing dashboard.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "uid": {
      "type": "string",
      "description": "Dashboard UID to update"
    },
    "dashboard": {
      "type": "object",
      "description": "Updated dashboard JSON configuration"
    }
  },
  "required": ["uid", "dashboard"]
}
```

**Output:** Success confirmation with version number.

### grafana_delete_dashboard

Delete a dashboard by UID.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "uid": {
      "type": "string",
      "description": "Dashboard UID to delete"
    }
  },
  "required": ["uid"]
}
```

**Output:** Success confirmation.

### grafana_search_dashboards

Search dashboards by query string or tags.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "query": {
      "type": "string",
      "description": "Search query"
    },
    "tags": {
      "type": "array",
      "items": {"type": "string"},
      "description": "Tags to filter by"
    }
  }
}
```

**Output:** Array of matching dashboards.

---

## Prometheus Tools

Prometheus provides metrics querying and time-series data analysis.

**Service:** Prometheus (ax-tools)
**URL:** `https://prometheus.axinova-internal.xyz`

### prometheus_query

Execute an instant Prometheus query.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "query": {
      "type": "string",
      "description": "PromQL query expression"
    },
    "time": {
      "type": "string",
      "description": "Evaluation timestamp (RFC3339 or Unix timestamp, optional)"
    }
  },
  "required": ["query"]
}
```

**Example:**
```bash
curl -X POST https://mcp.axinova-ai.com/api/mcp/v1/call \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "prometheus_query",
      "arguments": {
        "query": "up{job=\"node-exporter\"}"
      }
    }
  }'
```

**Output:** Query result with metric values and labels.

### prometheus_query_range

Execute a range Prometheus query over a time period.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "query": {
      "type": "string",
      "description": "PromQL query expression"
    },
    "start": {
      "type": "string",
      "description": "Start timestamp (RFC3339 or Unix timestamp)"
    },
    "end": {
      "type": "string",
      "description": "End timestamp (RFC3339 or Unix timestamp)"
    },
    "step": {
      "type": "string",
      "description": "Query resolution step width (e.g., '15s', '1m')"
    }
  },
  "required": ["query", "start", "end", "step"]
}
```

**Example:**
```bash
curl -X POST https://mcp.axinova-ai.com/api/mcp/v1/call \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "prometheus_query_range",
      "arguments": {
        "query": "rate(http_requests_total[5m])",
        "start": "2026-01-24T00:00:00Z",
        "end": "2026-01-24T23:59:59Z",
        "step": "1m"
      }
    }
  }'
```

**Output:** Time series data with values at each step interval.

### prometheus_series

Get list of time series that match label selectors.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "match": {
      "type": "array",
      "items": {"type": "string"},
      "description": "Series selectors (e.g., ['up', 'http_requests_total{job=\"api\"}'])"
    },
    "start": {
      "type": "string",
      "description": "Start timestamp"
    },
    "end": {
      "type": "string",
      "description": "End timestamp"
    }
  },
  "required": ["match"]
}
```

**Output:** Array of time series with their label sets.

### prometheus_labels

Get list of label names.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "start": {
      "type": "string",
      "description": "Start timestamp (optional)"
    },
    "end": {
      "type": "string",
      "description": "End timestamp (optional)"
    }
  }
}
```

**Output:** Array of label names.

### prometheus_label_values

Get list of label values for a specific label name.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "label": {
      "type": "string",
      "description": "Label name"
    }
  },
  "required": ["label"]
}
```

**Output:** Array of values for the specified label.

---

## SilverBullet Tools

SilverBullet provides wiki and knowledge base management.

**Service:** SilverBullet (ax-sas-tools)
**URL:** `https://wiki.axinova-internal.xyz`

### silverbullet_list_pages

List all wiki pages.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {}
}
```

**Example:**
```bash
curl -X POST https://mcp.axinova-ai.com/api/mcp/v1/call \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "silverbullet_list_pages",
      "arguments": {}
    }
  }'
```

**Output:** Array of page names and metadata.

### silverbullet_read_page

Read the content of a specific wiki page.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "page_name": {
      "type": "string",
      "description": "Name of the page to read"
    }
  },
  "required": ["page_name"]
}
```

**Example:**
```bash
curl -X POST https://mcp.axinova-ai.com/api/mcp/v1/call \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "silverbullet_read_page",
      "arguments": {"page_name": "Home"}
    }
  }'
```

**Output:** Page content in markdown format.

### silverbullet_write_page

Create or update a wiki page.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "page_name": {
      "type": "string",
      "description": "Name of the page"
    },
    "content": {
      "type": "string",
      "description": "Markdown content"
    }
  },
  "required": ["page_name", "content"]
}
```

**Output:** Success confirmation.

### silverbullet_delete_page

Delete a wiki page.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "page_name": {
      "type": "string",
      "description": "Name of the page to delete"
    }
  },
  "required": ["page_name"]
}
```

**Output:** Success confirmation.

### silverbullet_search_pages

Search for pages containing specific text.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "query": {
      "type": "string",
      "description": "Search query"
    }
  },
  "required": ["query"]
}
```

**Output:** Array of matching pages with excerpts.

---

## Vikunja Tools

Vikunja provides task and project management capabilities.

**Service:** Vikunja (ax-sas-tools)
**URL:** `https://vikunja.axinova-internal.xyz`

### vikunja_list_tasks

List all tasks, optionally filtered by project.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "project_id": {
      "type": "number",
      "description": "Filter by project ID (optional)"
    }
  }
}
```

**Example:**
```bash
curl -X POST https://mcp.axinova-ai.com/api/mcp/v1/call \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "vikunja_list_tasks",
      "arguments": {}
    }
  }'
```

**Output:** Array of task objects with ID, title, description, due date, status, etc.

### vikunja_get_task

Get details of a specific task.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "number",
      "description": "Task ID"
    }
  },
  "required": ["task_id"]
}
```

**Output:** Full task object with all details.

### vikunja_create_task

Create a new task.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "title": {
      "type": "string",
      "description": "Task title"
    },
    "description": {
      "type": "string",
      "description": "Task description (optional)"
    },
    "project_id": {
      "type": "number",
      "description": "Project ID to create task in"
    },
    "due_date": {
      "type": "string",
      "description": "Due date (RFC3339 format, optional)"
    },
    "priority": {
      "type": "number",
      "description": "Priority (1-5, optional)"
    }
  },
  "required": ["title", "project_id"]
}
```

**Example:**
```bash
curl -X POST https://mcp.axinova-ai.com/api/mcp/v1/call \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "vikunja_create_task",
      "arguments": {
        "title": "Review MCP documentation",
        "description": "Check all docs are up to date",
        "project_id": 1,
        "priority": 3
      }
    }
  }'
```

**Output:** Created task object with assigned ID.

### vikunja_update_task

Update an existing task.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "number",
      "description": "Task ID to update"
    },
    "title": {
      "type": "string",
      "description": "New title (optional)"
    },
    "description": {
      "type": "string",
      "description": "New description (optional)"
    },
    "done": {
      "type": "boolean",
      "description": "Mark as done/not done (optional)"
    },
    "due_date": {
      "type": "string",
      "description": "New due date (optional)"
    }
  },
  "required": ["task_id"]
}
```

**Output:** Updated task object.

### vikunja_delete_task

Delete a task.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "task_id": {
      "type": "number",
      "description": "Task ID to delete"
    }
  },
  "required": ["task_id"]
}
```

**Output:** Success confirmation.

### vikunja_list_projects

List all projects.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {}
}
```

**Output:** Array of project objects with ID, title, description, etc.

### vikunja_get_project

Get details of a specific project.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "project_id": {
      "type": "number",
      "description": "Project ID"
    }
  },
  "required": ["project_id"]
}
```

**Output:** Full project object with tasks count and metadata.

### vikunja_create_project

Create a new project.

**Input Schema:**
```json
{
  "type": "object",
  "properties": {
    "title": {
      "type": "string",
      "description": "Project title"
    },
    "description": {
      "type": "string",
      "description": "Project description (optional)"
    }
  },
  "required": ["title"]
}
```

**Output:** Created project object with assigned ID.

---

## Fetching Live Tool List

To get the most up-to-date list of tools with their exact schemas, query the API:

```bash
curl https://mcp.axinova-ai.com/api/mcp/v1/tools \
  -H "Authorization: Bearer YOUR_TOKEN" \
  | jq '.'
```

This will return the complete, current tool catalog with all input schemas and descriptions as registered in the running server.

## Tool Naming Convention

All tools follow the naming pattern: `{service}_{action}_{resource}`

**Examples:**
- `portainer_list_containers` - Portainer service, list action, containers resource
- `grafana_create_dashboard` - Grafana service, create action, dashboard resource
- `vikunja_update_task` - Vikunja service, update action, task resource

## Support

For questions about specific tools or to request new tools:
- GitHub Issues: https://github.com/axinova-ai/axinova-mcp-server-go/issues
- Documentation: https://github.com/axinova-ai/axinova-mcp-server-go/tree/main/docs
