# ChatGPT Integration Guide

**Platform:** OpenAI ChatGPT (GPT-4, GPT-3.5-turbo)
**Method:** Function Calling (Tools API)
**Last Updated:** 2026-01-25

> **💡 Note:** This guide is for integrating with **ChatGPT** (which doesn't support native MCP protocol).
>
> **For Claude Desktop, Claude Code, or GitHub Copilot**, use the much simpler **[Native MCP Integration](../QUICKSTART-NATIVE-MCP.md)** instead - it's 10x faster and requires no code!

---

## Overview

Integrate the Axinova MCP Server's 38 tools with ChatGPT using OpenAI's Function Calling feature. This allows ChatGPT to access your infrastructure tools (Portainer, Grafana, Prometheus, etc.) and execute operations based on natural language requests.

**This approach uses the HTTP API** because ChatGPT doesn't support the native MCP protocol.

---

## Prerequisites

- **OpenAI API Key** with GPT-4 or GPT-3.5-turbo access
- **MCP Server API Token**
- **Python 3.8+** and `openai` package

```bash
pip install openai requests
```

---

## Quick Start

### Step 1: Fetch MCP Tools

```python
import requests
import json

MCP_API_URL = "https://mcp.axinova-ai.com"
MCP_TOKEN = "your-mcp-token-here"

def get_mcp_tools():
    """Fetch all available MCP tools."""
    response = requests.get(
        f"{MCP_API_URL}/api/mcp/v1/tools",
        headers={"Authorization": f"Bearer {MCP_TOKEN}"}
    )
    response.raise_for_status()
    return response.json()["tools"]

mcp_tools = get_mcp_tools()
print(f"Loaded {len(mcp_tools)} tools")
```

### Step 2: Convert to OpenAI Function Format

```python
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

openai_functions = [convert_to_openai_function(t) for t in mcp_tools]
```

### Step 3: Create ChatGPT Conversation with Tools

```python
from openai import OpenAI

client = OpenAI(api_key="your-openai-api-key")

# Create conversation
messages = [
    {"role": "system", "content": "You are a DevOps assistant with access to infrastructure tools."},
    {"role": "user", "content": "List all Docker containers"}
]

response = client.chat.completions.create(
    model="gpt-4",
    messages=messages,
    tools=openai_functions,
    tool_choice="auto"
)

print(response.choices[0].message.content)
print(response.choices[0].message.tool_calls)
```

### Step 4: Execute Tool Calls

```python
def execute_mcp_tool(tool_name, arguments):
    """Execute an MCP tool and return the result."""
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
                "arguments": arguments
            }
        },
        timeout=30
    )

    response.raise_for_status()
    data = response.json()

    if "error" in data:
        raise Exception(f"MCP Error: {data['error']['message']}")

    return data["result"]["content"][0]["text"]

# Process tool calls
if response.choices[0].message.tool_calls:
    for tool_call in response.choices[0].message.tool_calls:
        function_name = tool_call.function.name
        function_args = json.loads(tool_call.function.arguments)

        print(f"Calling: {function_name}({function_args})")

        # Execute via MCP
        result = execute_mcp_tool(function_name, function_args)

        print(f"Result: {result}")
```

---

## Complete Integration Example

```python
#!/usr/bin/env python3
"""
ChatGPT + Axinova MCP Server Integration
Complete working example with function calling
"""

import requests
import json
from openai import OpenAI

# Configuration
MCP_API_URL = "https://mcp.axinova-ai.com"
MCP_TOKEN = "your-mcp-token-here"
OPENAI_API_KEY = "your-openai-key-here"

# Initialize OpenAI client
client = OpenAI(api_key=OPENAI_API_KEY)


def get_mcp_tools():
    """Fetch all MCP tools and convert to OpenAI format."""
    response = requests.get(
        f"{MCP_API_URL}/api/mcp/v1/tools",
        headers={"Authorization": f"Bearer {MCP_TOKEN}"}
    )
    response.raise_for_status()

    tools = response.json()["tools"]

    # Convert to OpenAI function format
    return [
        {
            "type": "function",
            "function": {
                "name": tool["name"],
                "description": tool["description"],
                "parameters": tool["inputSchema"]
            }
        }
        for tool in tools
    ]


def execute_mcp_tool(tool_name, arguments):
    """Execute MCP tool via HTTP API."""
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
                "arguments": arguments
            }
        },
        timeout=30
    )

    response.raise_for_status()
    data = response.json()

    if "error" in data:
        raise Exception(f"MCP Error: {data['error']['message']}")

    return data["result"]["content"][0]["text"]


def chat_with_tools(user_message, conversation_history=None):
    """Send a message to ChatGPT with MCP tools available."""
    if conversation_history is None:
        conversation_history = [
            {"role": "system", "content": "You are a DevOps assistant with access to infrastructure management tools including Docker (Portainer), monitoring (Grafana, Prometheus), wiki (SilverBullet), and tasks (Vikunja)."}
        ]

    # Add user message
    conversation_history.append({"role": "user", "content": user_message})

    # Get tools
    tools = get_mcp_tools()

    # Call ChatGPT
    response = client.chat.completions.create(
        model="gpt-4",
        messages=conversation_history,
        tools=tools,
        tool_choice="auto",
        temperature=0.7
    )

    message = response.choices[0].message

    # If ChatGPT wants to call tools
    if message.tool_calls:
        # Add assistant's tool call request to history
        conversation_history.append(message)

        # Execute each tool call
        for tool_call in message.tool_calls:
            function_name = tool_call.function.name
            function_args = json.loads(tool_call.function.arguments)

            print(f"\n🔧 Executing: {function_name}")
            print(f"   Arguments: {function_args}")

            try:
                # Execute tool
                result = execute_mcp_tool(function_name, function_args)

                # Add tool result to conversation
                conversation_history.append({
                    "role": "tool",
                    "tool_call_id": tool_call.id,
                    "name": function_name,
                    "content": result
                })

                print(f"   ✓ Success")

            except Exception as e:
                # Add error to conversation
                conversation_history.append({
                    "role": "tool",
                    "tool_call_id": tool_call.id,
                    "name": function_name,
                    "content": f"Error: {str(e)}"
                })

                print(f"   ✗ Error: {e}")

        # Get final response from ChatGPT with tool results
        final_response = client.chat.completions.create(
            model="gpt-4",
            messages=conversation_history,
            temperature=0.7
        )

        conversation_history.append(final_response.choices[0].message)
        return final_response.choices[0].message.content, conversation_history

    else:
        # No tool calls, just return the response
        conversation_history.append(message)
        return message.content, conversation_history


# Example usage
if __name__ == "__main__":
    print("ChatGPT + Axinova MCP Integration\n")

    # Single query
    response, history = chat_with_tools("List all Docker containers")
    print(f"\nChatGPT: {response}\n")

    # Multi-turn conversation
    response, history = chat_with_tools("Which containers are using the most memory?", history)
    print(f"\nChatGPT: {response}\n")

    response, history = chat_with_tools("Show me the Grafana dashboards related to containers", history)
    print(f"\nChatGPT: {response}\n")
```

---

## Example Conversations

### Example 1: Container Management

**User:** "List all Docker containers and show which ones are unhealthy"

**ChatGPT:** *Calls `portainer_list_containers`*

**Result:**
```
I found 12 containers running. Here's the health status:

✓ Healthy (10):
  - traefik (traefik:v3.0)
  - postgres (postgres:16)
  - grafana (grafana:latest)
  ...

⚠ Unhealthy (2):
  - redis-cache (redis:7) - Status: exited
  - worker-service (app:latest) - Status: restarting

Would you like me to get logs for the unhealthy containers?
```

### Example 2: Monitoring Queries

**User:** "What's the current CPU usage across all services?"

**ChatGPT:** *Calls `prometheus_query` with `rate(container_cpu_usage_seconds_total[5m])`*

**Result:**
```
Current CPU usage (5-minute average):
- traefik: 2.3%
- postgres: 15.7%
- grafana: 8.1%
- prometheus: 5.4%
- app: 45.2%  ← Highest

The app service is using the most CPU. Would you like me to investigate further?
```

### Example 3: Dashboard Analysis

**User:** "Show me all monitoring dashboards and tell me which ones are related to containers"

**ChatGPT:** *Calls `grafana_list_dashboards`*

**Result:**
```
I found 9 dashboards. Here are the container-related ones:

1. Container Resources
   - URL: /d/ffaprcnsaw2rkf/container-resources
   - Tags: containers, loki, prometheus

2. ECS Monitoring
   - URL: /d/afapto7uwpla8b/ecs-monitoring
   - Tags: host, infrastructure, mcp, registry

Would you like me to query specific metrics from any of these dashboards?
```

---

## Advanced Features

### 1. Conversation Memory

```python
# Maintain conversation across multiple queries
conversation_history = []

response1, conversation_history = chat_with_tools(
    "List all containers",
    conversation_history
)

response2, conversation_history = chat_with_tools(
    "Restart the ones that are unhealthy",  # Refers to previous context
    conversation_history
)
```

### 2. Streaming Responses

```python
response = client.chat.completions.create(
    model="gpt-4",
    messages=messages,
    tools=tools,
    stream=True
)

for chunk in response:
    if chunk.choices[0].delta.content:
        print(chunk.choices[0].delta.content, end="")
```

### 3. Error Handling

```python
def safe_execute_tool(tool_name, arguments):
    """Execute tool with comprehensive error handling."""
    try:
        return execute_mcp_tool(tool_name, arguments)
    except requests.Timeout:
        return f"Error: Tool execution timed out (30s)"
    except requests.HTTPError as e:
        if e.response.status_code == 401:
            return "Error: Invalid MCP API token"
        elif e.response.status_code == 429:
            return "Error: Rate limit exceeded, please retry later"
        else:
            return f"Error: HTTP {e.response.status_code}"
    except Exception as e:
        return f"Error: {str(e)}"
```

---

## Best Practices

### 1. Tool Selection

ChatGPT automatically selects appropriate tools, but you can guide it:

```python
# Encourage specific tool usage
messages = [
    {"role": "system", "content": "You have access to Portainer for Docker management, Grafana for dashboards, and Prometheus for metrics. Use these tools proactively to answer questions."},
    {"role": "user", "content": user_message}
]
```

### 2. Result Formatting

Process tool results before showing to user:

```python
# Parse JSON results for better display
result = execute_mcp_tool("portainer_list_containers", {"endpoint_id": 1})
containers = json.loads(result)

# Format as table
formatted = "| Name | Status | Image |\n"
formatted += "|------|--------|-------|\n"
for c in containers:
    formatted += f"| {c['Names']} | {c['State']} | {c['Image']} |\n"

# Add formatted result to conversation
conversation_history.append({
    "role": "tool",
    "tool_call_id": tool_call.id,
    "name": "portainer_list_containers",
    "content": formatted
})
```

### 3. Rate Limiting

```python
import time
from functools import wraps

def rate_limit(calls_per_minute=60):
    """Simple rate limiter decorator."""
    min_interval = 60.0 / calls_per_minute
    last_called = [0.0]

    def decorator(func):
        @wraps(func)
        def wrapper(*args, **kwargs):
            elapsed = time.time() - last_called[0]
            if elapsed < min_interval:
                time.sleep(min_interval - elapsed)
            result = func(*args, **kwargs)
            last_called[0] = time.time()
            return result
        return wrapper
    return decorator

@rate_limit(calls_per_minute=60)
def execute_mcp_tool(tool_name, arguments):
    # ... implementation
```

---

## Troubleshooting

### Issue: "Tool not found"

**Cause:** Tool name mismatch or typo

**Solution:**
```python
# Verify tool names
tools = get_mcp_tools()
tool_names = [t["function"]["name"] for t in tools]
print("Available tools:", tool_names)
```

### Issue: "Invalid arguments"

**Cause:** Argument type mismatch or missing required fields

**Solution:**
```python
# Check tool schema
tool = next(t for t in mcp_tools if t["name"] == "portainer_list_containers")
print("Required args:", tool["inputSchema"]["required"])
print("Properties:", tool["inputSchema"]["properties"])
```

### Issue: Rate limit errors

**Cause:** Too many requests in short time

**Solution:**
```python
# Add retry logic with exponential backoff
import time

def execute_with_retry(tool_name, arguments, max_retries=3):
    for attempt in range(max_retries):
        try:
            return execute_mcp_tool(tool_name, arguments)
        except requests.HTTPError as e:
            if e.response.status_code == 429 and attempt < max_retries - 1:
                wait_time = 2 ** attempt  # Exponential backoff
                print(f"Rate limited, waiting {wait_time}s...")
                time.sleep(wait_time)
            else:
                raise
```

---

## Next Steps

1. **Try the complete example** - Run the code and test with your own queries
2. **Customize system prompt** - Tailor ChatGPT's behavior for your use case
3. **Add error handling** - Implement robust error handling for production
4. **Build a UI** - Create a web interface with Streamlit or Gradio
5. **Monitor usage** - Track tool calls and costs

---

## See Also

- [Universal API Integration](../UNIVERSAL-API-INTEGRATION.md)
- [Gemini Integration](gemini.md)
- [LangChain Integration](langchain.md)
- [LlamaIndex Integration](llamaindex.md)

---

## Support

For issues or questions:
- GitHub Issues: https://github.com/axinova-ai/axinova-mcp-server-go/issues
- API Documentation: https://mcp.axinova-ai.com/docs
