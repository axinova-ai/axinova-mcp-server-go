#!/bin/bash

# Test script for MCP server

echo "Testing MCP Server..."
echo ""

# Start server
(
  # Send initialize request
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'

  # Send initialized notification
  echo '{"jsonrpc":"2.0","method":"initialized"}'

  # List tools
  echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}'

  # Wait a bit
  sleep 1
) | ./bin/axinova-mcp-server 2>&1
