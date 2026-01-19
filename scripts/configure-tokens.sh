#!/bin/bash

# Interactive token configuration script for MCP Server

set -e

SERVER="root@121.40.188.25"
DEPLOY_DIR="/opt/axinova-mcp-server"

echo "========================================="
echo "  MCP Server Token Configuration        "
echo "========================================="
echo ""
echo "This script will help you configure API tokens for the MCP server."
echo ""

# Check if running on server or remotely
if [ "$1" == "--remote" ]; then
    echo "📡 Configuring remotely via SSH..."
    ssh $SERVER "bash -s" < "$0" --local
    exit $?
fi

if [ "$1" != "--local" ]; then
    # Running from local machine
    echo "Run this script in one of two ways:"
    echo ""
    echo "1. Directly on server:"
    echo "   ssh $SERVER"
    echo "   cd $DEPLOY_DIR"
    echo "   bash scripts/configure-tokens.sh --local"
    echo ""
    echo "2. Remotely from local machine:"
    echo "   ./scripts/configure-tokens.sh --remote"
    echo ""
    exit 1
fi

# Running on server
cd $DEPLOY_DIR

echo "Current directory: $(pwd)"
echo ""

# Check if .env exists
if [ -f .env ]; then
    echo "⚠️  .env file already exists!"
    echo ""
    read -p "Backup existing .env? (y/n): " backup
    if [ "$backup" = "y" ]; then
        cp .env .env.backup.$(date +%Y%m%d%H%M%S)
        echo "✓ Backed up to .env.backup.$(date +%Y%m%d%H%M%S)"
    fi
    echo ""
fi

echo "========================================="
echo "  Step 1: Environment                   "
echo "========================================="
echo ""
echo "Set to 'prod' for production, 'dev' for development"
read -p "Environment (prod/dev) [prod]: " env_choice
ENV=${env_choice:-prod}

echo ""
echo "========================================="
echo "  Step 2: Portainer Token               "
echo "========================================="
echo ""
echo "How to get Portainer token:"
echo "  1. Go to https://portainer.axinova-internal.xyz"
echo "  2. Login → User settings → Access tokens"
echo "  3. Click 'Add access token'"
echo "  4. Name: mcp-server"
echo "  5. Copy the token (starts with ptr_)"
echo ""
read -p "Portainer token (or press Enter to skip): " portainer_token

echo ""
echo "========================================="
echo "  Step 3: Grafana Token                 "
echo "========================================="
echo ""
echo "How to get Grafana token:"
echo "  1. Go to https://grafana.axinova-internal.xyz"
echo "  2. Configuration (gear icon) → API Keys"
echo "  3. Click 'New API Key'"
echo "  4. Name: mcp-server, Role: Admin"
echo "  5. Copy the token (starts with glsa_ or eyJ)"
echo ""
read -p "Grafana token (or press Enter to skip): " grafana_token

echo ""
echo "========================================="
echo "  Step 4: SilverBullet Token            "
echo "========================================="
echo ""
echo "SilverBullet may not require authentication."
echo "Leave empty if you don't have a token."
echo ""
read -p "SilverBullet token (or press Enter to skip): " silverbullet_token

echo ""
echo "========================================="
echo "  Step 5: Vikunja Token                 "
echo "========================================="
echo ""
echo "How to get Vikunja token:"
echo "  1. Go to https://vikunja.axinova-internal.xyz"
echo "  2. Settings → API Tokens"
echo "  3. Click 'Create new token'"
echo "  4. Name: mcp-server"
echo "  5. Copy the token"
echo ""
read -p "Vikunja token (or press Enter to skip): " vikunja_token

# Create .env file
echo ""
echo "========================================="
echo "  Creating .env file...                 "
echo "========================================="
echo ""

cat > .env << ENVEOF
# Environment
ENV=$ENV

# ============================================
# Portainer (ax-tools)
# ============================================
APP_PORTAINER__URL=https://portainer.axinova-internal.xyz
APP_PORTAINER__TOKEN=$portainer_token

# ============================================
# Grafana (ax-tools)
# ============================================
APP_GRAFANA__URL=https://grafana.axinova-internal.xyz
APP_GRAFANA__TOKEN=$grafana_token

# ============================================
# Prometheus (ax-tools)
# ============================================
APP_PROMETHEUS__URL=https://prometheus.axinova-internal.xyz

# ============================================
# SilverBullet (ax-sas-tools)
# ============================================
APP_SILVERBULLET__URL=https://silverbullet.axinova-internal.xyz
APP_SILVERBULLET__TOKEN=$silverbullet_token

# ============================================
# Vikunja (ax-sas-tools)
# ============================================
APP_VIKUNJA__URL=https://vikunja.axinova-internal.xyz
APP_VIKUNJA__TOKEN=$vikunja_token
ENVEOF

# Set secure permissions
chmod 600 .env

echo "✓ .env file created"
echo "✓ Permissions set to 600 (owner read/write only)"
echo ""

# Summary
echo "========================================="
echo "  Configuration Summary                 "
echo "========================================="
echo ""
echo "Environment: $ENV"
echo ""
echo "Configured services:"
[ -n "$portainer_token" ] && echo "  ✓ Portainer" || echo "  ⊗ Portainer (no token)"
[ -n "$grafana_token" ] && echo "  ✓ Grafana" || echo "  ⊗ Grafana (no token)"
echo "  ✓ Prometheus (no token required)"
[ -n "$silverbullet_token" ] && echo "  ✓ SilverBullet" || echo "  ⊗ SilverBullet (no token)"
[ -n "$vikunja_token" ] && echo "  ✓ Vikunja" || echo "  ⊗ Vikunja (no token)"
echo ""

echo "========================================="
echo "  Next Steps                            "
echo "========================================="
echo ""
echo "1. Start the service:"
echo "   docker-compose up -d"
echo ""
echo "2. Check logs:"
echo "   docker-compose logs -f mcp-server"
echo ""
echo "3. Test the server:"
echo "   docker-compose exec mcp-server /app/axinova-mcp-server"
echo ""
echo "For more help, see:"
echo "  - DEPLOYMENT.md"
echo "  - TESTING.md"
echo "  - scripts/get-tokens.md"
echo ""
