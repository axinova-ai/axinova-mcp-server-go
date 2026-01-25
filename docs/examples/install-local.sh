#!/bin/bash
# Install Axinova MCP Server locally
#
# Usage:
#   ./install-local.sh
#
# This script will:
# - Detect your OS (macOS or Linux)
# - Download the latest release from GitHub
# - Install to /usr/local/bin/axinova-mcp-server
# - Make it executable

set -e

echo "Installing Axinova MCP Server..."
echo ""

# Detect OS
OS=$(uname -s)
ARCH=$(uname -m)

case "$OS" in
  Darwin)
    BINARY="axinova-mcp-server-macos"
    echo "Detected: macOS ($ARCH)"
    ;;
  Linux)
    BINARY="axinova-mcp-server-linux"
    echo "Detected: Linux ($ARCH)"
    ;;
  *)
    echo "❌ Unsupported OS: $OS"
    echo "   Supported: macOS (Darwin), Linux"
    exit 1
    ;;
esac

echo ""
echo "Downloading latest release..."

# Download latest release
GITHUB_REPO="axinova-ai/axinova-mcp-server-go"
DOWNLOAD_URL="https://github.com/$GITHUB_REPO/releases/latest/download/$BINARY"

curl -L "$DOWNLOAD_URL" -o /tmp/axinova-mcp-server

if [ ! -f /tmp/axinova-mcp-server ]; then
  echo "❌ Download failed"
  echo "   URL: $DOWNLOAD_URL"
  exit 1
fi

echo "✓ Downloaded to /tmp/axinova-mcp-server"
echo ""

# Install to /usr/local/bin
echo "Installing to /usr/local/bin/axinova-mcp-server..."
echo "(You may be prompted for your password)"

sudo mv /tmp/axinova-mcp-server /usr/local/bin/axinova-mcp-server
sudo chmod +x /usr/local/bin/axinova-mcp-server

echo "✓ Installed"
echo ""

# Verify installation
echo "Verifying installation..."

if command -v axinova-mcp-server &> /dev/null; then
  echo "✅ Installation successful!"
  echo ""
  echo "Location: $(which axinova-mcp-server)"
else
  echo "⚠️  Installation complete, but binary not in PATH"
  echo "   This may be normal. Try:"
  echo "   /usr/local/bin/axinova-mcp-server --help"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Next steps:"
echo ""
echo "1. Configure your MCP client:"
echo "   - Claude Desktop: Edit ~/Library/Application Support/Claude/claude_desktop_config.json"
echo "   - Claude Code: Run 'claude mcp add axinova-tools --scope user -- /usr/local/bin/axinova-mcp-server'"
echo "   - GitHub Copilot: Add to VS Code settings (see docs)"
echo ""
echo "2. See documentation:"
echo "   - Main guide: docs/NATIVE-MCP-INTEGRATION.md"
echo "   - Claude Desktop: docs/onboarding/claude-desktop.md"
echo "   - Claude Code: docs/onboarding/claude-code.md"
echo "   - GitHub Copilot: docs/onboarding/github-copilot.md"
echo ""
echo "3. Example configurations:"
echo "   - docs/examples/claude_desktop_config.json"
echo "   - docs/examples/vscode_settings.json"
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
