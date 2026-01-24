# Integrating MCP Server with LLM Agents

## Overview

The MCP (Model Context Protocol) server is designed to be consumed by LLM agents and assistants. This guide shows how to integrate it with various LLM platforms and frameworks.

The server provides two integration modes:
1. **Native MCP Protocol (stdio)** - For MCP-compliant clients like Claude Desktop
2. **HTTP JSON-RPC API** - For custom agents, LangChain, LlamaIndex, and other frameworks

## Quick Start

```bash
# Test the API is accessible
curl https://mcp.axinova-ai.com/health

# List available tools
curl https://mcp.axinova-ai.com/api/mcp/v1/tools \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"

# Call a tool
curl -X POST https://mcp.axinova-ai.com/api/mcp/v1/call \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
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

## Integration Methods

### 1. Claude Desktop (MCP Native)

Claude Desktop supports the MCP protocol natively via two modes: stdio (local) and SSE (remote).

#### stdio Mode (Local Development)

For running the MCP server locally on your machine:

**Configure `claude_desktop_config.json`:**

**macOS:** `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows:** `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "axinova-tools": {
      "command": "/path/to/axinova-mcp-server-go",
      "args": [],
      "env": {
        "ENV": "prod",
        "APP_PORTAINER__URL": "https://portainer.axinova-internal.xyz",
        "APP_PORTAINER__TOKEN": "ptr_...",
        "APP_PORTAINER__ENABLED": "true",
        "APP_GRAFANA__URL": "https://grafana.axinova-internal.xyz",
        "APP_GRAFANA__TOKEN": "...",
        "APP_GRAFANA__ENABLED": "true",
        "APP_PROMETHEUS__URL": "https://prometheus.axinova-internal.xyz",
        "APP_PROMETHEUS__ENABLED": "true"
      }
    }
  }
}
```

After updating the config:
1. Restart Claude Desktop
2. The MCP server tools will appear in Claude's tool palette
3. Claude can now use tools like `portainer_list_containers`, `grafana_list_dashboards`, etc.

#### HTTP Mode (Remote Server)

For using the deployed MCP server via HTTP, you need to bridge the HTTP API to Claude Desktop. Options:

**Option A: Use the Claude API with Custom Tools**

Convert MCP tools to Claude API tool format and call the HTTP API in your application:

```python
import anthropic
import requests

# Fetch MCP tools
mcp_tools = requests.get(
    "https://mcp.axinova-ai.com/api/mcp/v1/tools",
    headers={"Authorization": "Bearer YOUR_TOKEN"}
).json()["tools"]

# Convert to Claude API format
claude_tools = []
for tool in mcp_tools:
    claude_tools.append({
        "name": tool["name"],
        "description": tool["description"],
        "input_schema": tool["inputSchema"]
    })

# Use with Claude API
client = anthropic.Anthropic(api_key="YOUR_ANTHROPIC_KEY")

message = client.messages.create(
    model="claude-3-5-sonnet-20241022",
    max_tokens=1024,
    tools=claude_tools,
    messages=[{"role": "user", "content": "List Docker containers"}]
)

# When Claude requests a tool use, call the MCP API
if message.stop_reason == "tool_use":
    tool_use = message.content[1]  # Get tool use block

    result = requests.post(
        "https://mcp.axinova-ai.com/api/mcp/v1/call",
        headers={"Authorization": "Bearer YOUR_TOKEN"},
        json={
            "jsonrpc": "2.0",
            "id": 1,
            "method": "tools/call",
            "params": {
                "name": tool_use.name,
                "arguments": tool_use.input
            }
        }
    ).json()["result"]

    # Send result back to Claude
    # ... continue conversation
```

**Option B: MCP-over-HTTP Proxy**

Create a local proxy that implements MCP stdio and forwards to the HTTP API:

```python
# mcp_http_proxy.py
import sys
import json
import requests

API_URL = "https://mcp.axinova-ai.com"
API_TOKEN = "sk-mcp-prod-..."

def handle_request(request):
    if request["method"] == "initialize":
        return {
            "protocolVersion": "2025-11-25",
            "capabilities": {"tools": {}},
            "serverInfo": {"name": "axinova-mcp-http-proxy", "version": "1.0.0"}
        }

    elif request["method"] == "tools/list":
        response = requests.get(
            f"{API_URL}/api/mcp/v1/tools",
            headers={"Authorization": f"Bearer {API_TOKEN}"}
        )
        return {"tools": response.json()["tools"]}

    elif request["method"] == "tools/call":
        response = requests.post(
            f"{API_URL}/api/mcp/v1/call",
            headers={"Authorization": f"Bearer {API_TOKEN}"},
            json=request
        )
        return response.json()["result"]

# MCP stdio loop
for line in sys.stdin:
    request = json.loads(line)
    result = handle_request(request)
    response = {"jsonrpc": "2.0", "id": request["id"], "result": result}
    print(json.dumps(response), flush=True)
```

Then configure Claude Desktop:
```json
{
  "mcpServers": {
    "axinova-tools-remote": {
      "command": "python3",
      "args": ["/path/to/mcp_http_proxy.py"]
    }
  }
}
```

### 2. LangChain Integration

LangChain can use the MCP server tools via custom tool wrappers.

```python
from langchain.tools import Tool
from langchain.agents import initialize_agent, AgentType
from langchain.llms import OpenAI
import requests

MCP_API_URL = "https://mcp.axinova-ai.com"
MCP_API_TOKEN = "sk-mcp-prod-..."

def call_mcp_tool(tool_name: str, arguments: dict) -> str:
    """Call an MCP tool and return the result"""
    response = requests.post(
        f"{MCP_API_URL}/api/mcp/v1/call",
        headers={"Authorization": f"Bearer {MCP_API_TOKEN}"},
        json={
            "jsonrpc": "2.0",
            "id": 1,
            "method": "tools/call",
            "params": {
                "name": tool_name,
                "arguments": arguments
            }
        }
    )

    data = response.json()
    if "error" in data:
        return f"Error: {data['error']['message']}"

    # Extract text content from MCP result
    content = data["result"]["content"]
    return content[0]["text"] if content else ""

# Fetch available tools
tools_response = requests.get(
    f"{MCP_API_URL}/api/mcp/v1/tools",
    headers={"Authorization": f"Bearer {MCP_API_TOKEN}"}
).json()

# Convert MCP tools to LangChain tools
langchain_tools = []
for tool in tools_response["tools"]:
    def make_tool_func(tool_name):
        return lambda args: call_mcp_tool(tool_name, eval(args) if isinstance(args, str) else args)

    langchain_tools.append(Tool(
        name=tool["name"],
        description=tool["description"],
        func=make_tool_func(tool["name"])
    ))

# Create agent
llm = OpenAI(temperature=0)
agent = initialize_agent(
    langchain_tools,
    llm,
    agent=AgentType.ZERO_SHOT_REACT_DESCRIPTION,
    verbose=True
)

# Use the agent
result = agent.run("List all Docker containers in Portainer endpoint 1")
print(result)
```

### 3. LlamaIndex Integration

LlamaIndex can use MCP tools as custom function tools:

```python
from llama_index.core.tools import FunctionTool
from llama_index.core.agent import ReActAgent
from llama_index.llms.openai import OpenAI
import requests

MCP_API_URL = "https://mcp.axinova-ai.com"
MCP_API_TOKEN = "sk-mcp-prod-..."

def create_mcp_tool(tool_name: str, tool_description: str, tool_schema: dict):
    """Create a LlamaIndex FunctionTool that calls the MCP API"""

    def tool_function(**kwargs):
        response = requests.post(
            f"{MCP_API_URL}/api/mcp/v1/call",
            headers={"Authorization": f"Bearer {MCP_API_TOKEN}"},
            json={
                "jsonrpc": "2.0",
                "id": 1,
                "method": "tools/call",
                "params": {
                    "name": tool_name,
                    "arguments": kwargs
                }
            }
        )

        data = response.json()
        if "error" in data:
            return f"Error: {data['error']['message']}"

        content = data["result"]["content"]
        return content[0]["text"] if content else ""

    return FunctionTool.from_defaults(
        fn=tool_function,
        name=tool_name,
        description=tool_description
    )

# Fetch and convert MCP tools
tools_response = requests.get(
    f"{MCP_API_URL}/api/mcp/v1/tools",
    headers={"Authorization": f"Bearer {MCP_API_TOKEN}"}
).json()

llamaindex_tools = []
for tool in tools_response["tools"]:
    llamaindex_tools.append(
        create_mcp_tool(tool["name"], tool["description"], tool["inputSchema"])
    )

# Create agent
llm = OpenAI(model="gpt-4")
agent = ReActAgent.from_tools(llamaindex_tools, llm=llm, verbose=True)

# Use the agent
response = agent.chat("List all Grafana dashboards")
print(response)
```

### 4. OpenAI Function Calling

Convert MCP tools to OpenAI function schemas:

```python
import openai
import requests

MCP_API_URL = "https://mcp.axinova-ai.com"
MCP_API_TOKEN = "sk-mcp-prod-..."

# Fetch MCP tools
tools_response = requests.get(
    f"{MCP_API_URL}/api/mcp/v1/tools",
    headers={"Authorization": f"Bearer {MCP_API_TOKEN}"}
).json()

# Convert to OpenAI function format
openai_functions = []
for tool in tools_response["tools"]:
    openai_functions.append({
        "type": "function",
        "function": {
            "name": tool["name"],
            "description": tool["description"],
            "parameters": tool["inputSchema"]
        }
    })

# Use with OpenAI API
client = openai.OpenAI()

messages = [{"role": "user", "content": "List all Docker containers"}]

response = client.chat.completions.create(
    model="gpt-4",
    messages=messages,
    tools=openai_functions,
    tool_choice="auto"
)

# Handle tool calls
if response.choices[0].message.tool_calls:
    for tool_call in response.choices[0].message.tool_calls:
        # Call MCP API
        mcp_response = requests.post(
            f"{MCP_API_URL}/api/mcp/v1/call",
            headers={"Authorization": f"Bearer {MCP_API_TOKEN}"},
            json={
                "jsonrpc": "2.0",
                "id": 1,
                "method": "tools/call",
                "params": {
                    "name": tool_call.function.name,
                    "arguments": eval(tool_call.function.arguments)
                }
            }
        ).json()

        result = mcp_response["result"]["content"][0]["text"]

        # Add tool result to messages
        messages.append({
            "role": "tool",
            "tool_call_id": tool_call.id,
            "content": result
        })

    # Get final response
    final_response = client.chat.completions.create(
        model="gpt-4",
        messages=messages
    )
    print(final_response.choices[0].message.content)
```

### 5. Custom LLM Agent

Generic HTTP API integration for any LLM framework:

```python
import requests
from typing import Dict, List, Any

class MCPClient:
    def __init__(self, api_url: str, api_token: str):
        self.api_url = api_url
        self.headers = {
            "Authorization": f"Bearer {api_token}",
            "Content-Type": "application/json"
        }

    def list_tools(self) -> List[Dict[str, Any]]:
        """Fetch all available MCP tools"""
        response = requests.get(
            f"{self.api_url}/api/mcp/v1/tools",
            headers=self.headers
        )
        response.raise_for_status()
        return response.json()["tools"]

    def call_tool(self, tool_name: str, arguments: Dict[str, Any]) -> str:
        """Execute an MCP tool and return the result"""
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

        # Extract text content
        content = data["result"]["content"]
        return content[0]["text"] if content else ""

    def get_tool_schema(self, tool_name: str) -> Dict[str, Any]:
        """Get the input schema for a specific tool"""
        tools = self.list_tools()
        for tool in tools:
            if tool["name"] == tool_name:
                return tool["inputSchema"]
        raise ValueError(f"Tool '{tool_name}' not found")

# Usage in your LLM agent
mcp = MCPClient(
    "https://mcp.axinova-ai.com",
    "sk-mcp-prod-86d850ac73a8b9dd11e94b104ea4fd56966bee365ed5ffa3820ecd99f5f2640e"
)

# Discover available tools
tools = mcp.list_tools()
print(f"Available tools: {len(tools)}")

# Call a tool
containers = mcp.call_tool("portainer_list_containers", {"endpoint_id": 1})
print(containers)

# Get tool schema for validation
schema = mcp.get_tool_schema("grafana_list_dashboards")
print(f"Tool schema: {schema}")
```

## MCP Protocol Compatibility

### Native MCP Support (stdio)

The server implements the MCP protocol specification (2025-11-25) for stdio transport:

1. **Initialize handshake**: Client sends `initialize` request with protocol version
2. **Capabilities negotiation**: Server responds with supported features
3. **Tool discovery**: Client calls `tools/list` method
4. **Tool execution**: Client calls `tools/call` method
5. **Resource access**: Client calls `resources/read` method (if supported)

**Compatible with:**
- Claude Desktop
- Any MCP-compliant client implementing the 2025-11-25 specification

**Example stdio session:**
```json
// Client sends
{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2025-11-25"}}

// Server responds
{"jsonrpc": "2.0", "id": 1, "result": {"protocolVersion": "2025-11-25", "capabilities": {"tools": {}}}}

// Client lists tools
{"jsonrpc": "2.0", "id": 2, "method": "tools/list"}

// Server responds with tools array
{"jsonrpc": "2.0", "id": 2, "result": {"tools": [...]}}
```

### HTTP JSON-RPC Compatibility

The HTTP API provides the same MCP methods over HTTP, compatible with:
- LangChain
- LlamaIndex
- AutoGPT
- Custom agents
- Any HTTP client

**Transport differences:**
- **stdio**: Newline-delimited JSON over stdin/stdout
- **HTTP**: JSON-RPC over HTTPS with Bearer token authentication

## Security Considerations

### Token Management

1. **Storage**: Store API tokens securely
   - Use environment variables
   - Use secrets managers (AWS Secrets Manager, HashiCorp Vault)
   - Never commit tokens to git repositories

2. **Rotation**: Rotate tokens periodically (recommended: every 90 days)

3. **Scope**: Each token has access to all enabled services
   - Consider creating separate MCP server instances for different scopes if needed

### Network Security

1. **Transport**: All HTTP communication uses HTTPS (TLS 1.3)
2. **Internal services**: Backing services (Portainer, Grafana) are on private network
3. **Firewall**: MCP server is the only public-facing endpoint
4. **Rate limiting**: 1000 requests per minute (configurable)

### Error Handling

Never expose internal service details in errors:
```python
try:
    result = mcp.call_tool("portainer_list_containers", {"endpoint_id": 1})
except Exception as e:
    # Log full error internally
    logger.error(f"MCP call failed: {e}")

    # Return generic error to user
    return "Failed to fetch containers. Please try again."
```

## Performance Optimization

### Caching

Cache tool results when appropriate:
```python
import functools
import time

@functools.lru_cache(maxsize=128)
def cached_mcp_call(tool_name: str, args_hash: str) -> str:
    # Reconstruct args from hash
    return mcp.call_tool(tool_name, eval(args_hash))

# Use with TTL
def call_with_ttl(tool_name: str, arguments: dict, ttl_seconds: int = 300):
    cache_key = (tool_name, str(sorted(arguments.items())), time.time() // ttl_seconds)
    return cached_mcp_call(cache_key[0], cache_key[1])
```

### Parallel Requests

Make parallel tool calls when possible:
```python
import asyncio
import aiohttp

async def call_tool_async(session, tool_name, arguments):
    async with session.post(
        f"{MCP_API_URL}/api/mcp/v1/call",
        headers={"Authorization": f"Bearer {MCP_API_TOKEN}"},
        json={
            "jsonrpc": "2.0",
            "id": 1,
            "method": "tools/call",
            "params": {"name": tool_name, "arguments": arguments}
        }
    ) as response:
        return await response.json()

async def parallel_calls():
    async with aiohttp.ClientSession() as session:
        tasks = [
            call_tool_async(session, "portainer_list_containers", {"endpoint_id": 1}),
            call_tool_async(session, "grafana_list_dashboards", {}),
            call_tool_async(session, "vikunja_list_tasks", {})
        ]
        results = await asyncio.gather(*tasks)
        return results

# Run
results = asyncio.run(parallel_calls())
```

## Example: Complete Agent Implementation

Here's a complete example of a simple CLI agent using the MCP server:

```python
#!/usr/bin/env python3
import requests
import json
from typing import Dict, Any

class SimpleAgent:
    def __init__(self, mcp_url: str, mcp_token: str):
        self.mcp_url = mcp_url
        self.mcp_token = mcp_token
        self.tools = self._load_tools()

    def _load_tools(self):
        """Fetch available tools"""
        response = requests.get(
            f"{self.mcp_url}/api/mcp/v1/tools",
            headers={"Authorization": f"Bearer {self.mcp_token}"}
        )
        return {tool["name"]: tool for tool in response.json()["tools"]}

    def call_tool(self, tool_name: str, arguments: Dict[str, Any]) -> str:
        """Execute a tool"""
        payload = {
            "jsonrpc": "2.0",
            "id": 1,
            "method": "tools/call",
            "params": {"name": tool_name, "arguments": arguments}
        }

        response = requests.post(
            f"{self.mcp_url}/api/mcp/v1/call",
            headers={"Authorization": f"Bearer {self.mcp_token}"},
            json=payload
        )

        data = response.json()
        if "error" in data:
            return f"Error: {data['error']['message']}"

        return data["result"]["content"][0]["text"]

    def run(self):
        """Simple REPL"""
        print(f"Agent initialized with {len(self.tools)} tools")
        print("Available tools:", ", ".join(self.tools.keys()))
        print("\nType 'list' to see tools, 'exit' to quit\n")

        while True:
            command = input("> ").strip()

            if command == "exit":
                break
            elif command == "list":
                for name, tool in self.tools.items():
                    print(f"  {name}: {tool['description']}")
            else:
                # Simple command parsing (tool_name arg1=val1 arg2=val2)
                parts = command.split()
                if not parts:
                    continue

                tool_name = parts[0]
                if tool_name not in self.tools:
                    print(f"Unknown tool: {tool_name}")
                    continue

                # Parse arguments
                args = {}
                for part in parts[1:]:
                    if "=" in part:
                        key, value = part.split("=", 1)
                        # Try to parse as JSON, fall back to string
                        try:
                            args[key] = json.loads(value)
                        except:
                            args[key] = value

                # Call tool
                result = self.call_tool(tool_name, args)
                print(result)

if __name__ == "__main__":
    agent = SimpleAgent(
        "https://mcp.axinova-ai.com",
        "sk-mcp-prod-86d850ac73a8b9dd11e94b104ea4fd56966bee365ed5ffa3820ecd99f5f2640e"
    )
    agent.run()
```

## Troubleshooting

### Common Issues

**Issue: "Missing Authorization header"**
- Ensure you're sending the `Authorization: Bearer <token>` header
- Verify the token is correct (check for extra spaces or newlines)

**Issue: "Service unavailable"**
- The backing service (Portainer, Grafana, etc.) may be down
- Check service health: `curl https://portainer.axinova-internal.xyz`
- Check MCP server logs for more details

**Issue: "Invalid params"**
- Verify your arguments match the tool's input schema
- Use `GET /api/mcp/v1/tools` to see required parameters
- Check JSON formatting in your request

**Issue: Rate limit exceeded**
- Wait for the rate limit window to reset
- Implement exponential backoff in your client
- Contact support to increase limits if needed

### Debugging

Enable verbose logging in your client:
```python
import logging

logging.basicConfig(level=logging.DEBUG)

# Now all requests will be logged
result = mcp.call_tool("portainer_list_containers", {"endpoint_id": 1})
```

Check MCP server metrics:
```bash
curl https://mcp.axinova-ai.com/metrics | grep mcp_rpc_errors_total
```

## Next Steps

1. Read the [API Reference](./API-REFERENCE.md) for detailed endpoint documentation
2. Browse the [Tool Catalog](./TOOL-CATALOG.md) to see all available tools
3. Try the example agents above with your own use cases
4. Build custom integrations for your specific LLM framework

## Support

For questions or issues:
- GitHub Issues: https://github.com/axinova-ai/axinova-mcp-server-go/issues
- Documentation: https://github.com/axinova-ai/axinova-mcp-server-go/tree/main/docs
