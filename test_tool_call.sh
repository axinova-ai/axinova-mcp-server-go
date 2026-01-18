#!/bin/bash

# Test calling an actual tool (Prometheus query)

echo "Testing MCP Tool Call (Prometheus)..."
echo ""

(
  # Initialize
  echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test-client","version":"1.0.0"}}}'

  # Initialized
  echo '{"jsonrpc":"2.0","method":"initialized"}'

  # Call prometheus_query tool
  echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"prometheus_query","arguments":{"query":"up"}}}'

  # Wait
  sleep 2
) | ./bin/axinova-mcp-server 2>&1 | grep -A 10 '"id":3'
