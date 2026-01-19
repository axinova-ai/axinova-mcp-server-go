#!/bin/bash

# Deployment script for ax-sas-tools (121.40.188.25)

set -e

SERVER="root@121.40.188.25"
DEPLOY_DIR="/opt/axinova-mcp-server"
REPO_URL="https://github.com/axinova-ai/axinova-mcp-server-go.git"

echo "========================================="
echo "  Deploying MCP Server to ax-sas-tools  "
echo "========================================="
echo ""

# Check SSH connection
echo "🔍 Checking SSH connection..."
ssh -o ConnectTimeout=5 $SERVER "echo '✓ SSH connection OK'" || {
    echo "❌ Failed to connect to $SERVER"
    exit 1
}

# Deploy
echo ""
echo "📦 Deploying to $SERVER:$DEPLOY_DIR..."

ssh $SERVER << 'ENDSSH'
set -e

# Variables
DEPLOY_DIR="/opt/axinova-mcp-server"
REPO_URL="https://github.com/axinova-ai/axinova-mcp-server-go.git"

echo "📁 Creating deployment directory..."
mkdir -p $DEPLOY_DIR
cd $DEPLOY_DIR

# Clone or update repository
if [ -d ".git" ]; then
    echo "📥 Updating repository..."
    git fetch origin
    git reset --hard origin/main
    git pull origin main
else
    echo "📥 Cloning repository..."
    git clone $REPO_URL .
fi

echo "✓ Code deployed"

# Build Docker image
echo ""
echo "🐳 Building Docker image..."
docker build -t axinova-mcp-server:latest .

echo "✓ Docker image built"

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo ""
    echo "⚠️  Creating .env file from example..."
    cp .env.example .env
    echo ""
    echo "⚠️  IMPORTANT: Edit /opt/axinova-mcp-server/.env with actual API tokens!"
    echo "              See scripts/get-tokens.md for token generation instructions"
else
    echo "✓ .env file exists"
fi

# Show deployment status
echo ""
echo "========================================="
echo "  Deployment Complete!                  "
echo "========================================="
echo ""
echo "Location: $DEPLOY_DIR"
echo "Image: axinova-mcp-server:latest"
echo ""
echo "Next steps:"
echo "  1. Configure API tokens:"
echo "     Edit /opt/axinova-mcp-server/.env"
echo ""
echo "  2. Start the service:"
echo "     cd /opt/axinova-mcp-server"
echo "     docker-compose up -d"
echo ""
echo "  3. Check logs:"
echo "     docker-compose logs -f mcp-server"
echo ""
echo "  4. Test the server:"
echo "     docker-compose exec mcp-server /app/axinova-mcp-server"
echo ""

ENDSSH

echo ""
echo "✅ Deployment script completed successfully!"
echo ""
echo "To configure and start:"
echo "  ssh $SERVER"
echo "  cd $DEPLOY_DIR"
echo "  # Edit .env with tokens"
echo "  docker-compose up -d"
