# Deployment Guide

This document provides step-by-step instructions for deploying the Axinova MCP Server to production.

## Current Deployment Status

✅ **Deployed to:** ax-sas-tools (121.40.188.25)
✅ **Location:** `/opt/axinova-mcp-server`
✅ **Image:** `axinova-mcp-server:latest`
✅ **Build:** Successful

## Prerequisites

- SSH access to deployment server
- API tokens for all services
- Docker and docker-compose installed on server

## Quick Start

### 1. Configure API Tokens

SSH to the server and edit the environment file:

```bash
ssh root@121.40.188.25
cd /opt/axinova-mcp-server
nano .env
```

Add your API tokens (see [scripts/get-tokens.md](scripts/get-tokens.md) for token generation):

```bash
# Portainer
APP_PORTAINER__TOKEN=ptr_your_portainer_token

# Grafana
APP_GRAFANA__TOKEN=glsa_your_grafana_token

# SilverBullet (optional, leave empty if no auth)
APP_SILVERBULLET__TOKEN=

# Vikunja
APP_VIKUNJA__TOKEN=tk_your_vikunja_token
```

### 2. Start the Service

```bash
docker-compose up -d
```

### 3. Verify Deployment

```bash
# Check container status
docker-compose ps

# View logs
docker-compose logs -f mcp-server

# Test the server
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}' | docker-compose exec -T mcp-server /app/axinova-mcp-server
```

## Deployment Architecture

### Server Locations

| Service | Host | URL |
|---------|------|-----|
| Portainer | ax-tools | https://portainer.axinova-internal.xyz |
| Grafana | ax-tools | https://grafana.axinova-internal.xyz |
| Prometheus | ax-tools | https://prometheus.axinova-internal.xyz |
| SilverBullet | ax-sas-tools | https://silverbullet.axinova-internal.xyz |
| Vikunja | ax-sas-tools | https://vikunja.axinova-internal.xyz |
| **MCP Server** | **ax-sas-tools** | **(121.40.188.25)** |

### Network Configuration

- TLS termination: Traefik
- TLS verification: Disabled (`tls.skip_verify: true` in prod config)
- All services accessible via HTTPS with internal certificates

## Automated Deployment

Use the deployment script for automated deployment:

```bash
# From your local machine
cd axinova-mcp-server-go
./scripts/deploy-to-ax-sas-tools.sh
```

The script will:
1. Check SSH connection
2. Clone/update repository
3. Build Docker image
4. Create .env file if needed

## Manual Deployment

### Step 1: SSH to Server

```bash
ssh root@121.40.188.25
```

### Step 2: Clone Repository

```bash
mkdir -p /opt/axinova-mcp-server
cd /opt/axinova-mcp-server
git clone https://github.com/axinova-ai/axinova-mcp-server-go.git .
```

### Step 3: Build Docker Image

```bash
docker build -t axinova-mcp-server:latest .
```

### Step 4: Configure Environment

```bash
cp .env.example .env
nano .env
# Add your API tokens
```

### Step 5: Start with Docker Compose

```bash
docker-compose up -d
```

## Configuration

### Environment Variables

All configuration can be set via environment variables with `APP_` prefix:

```bash
# Service URLs (already set in config/prod.yaml)
APP_PORTAINER__URL=https://portainer.axinova-internal.xyz
APP_GRAFANA__URL=https://grafana.axinova-internal.xyz
APP_PROMETHEUS__URL=https://prometheus.axinova-internal.xyz
APP_SILVERBULLET__URL=https://silverbullet.axinova-internal.xyz
APP_VIKUNJA__URL=https://vikunja.axinova-internal.xyz

# API Tokens (REQUIRED - set in .env)
APP_PORTAINER__TOKEN=your-token
APP_GRAFANA__TOKEN=your-token
APP_SILVERBULLET__TOKEN=your-token
APP_VIKUNJA__TOKEN=your-token

# TLS (already set in config/prod.yaml)
APP_TLS__SKIP_VERIFY=true

# Optional: Disable specific services
APP_PORTAINER__ENABLED=false
APP_GRAFANA__ENABLED=false
```

### Configuration Files

- `config/base.yaml` - Base configuration
- `config/prod.yaml` - Production overrides
- `.env` - Environment-specific secrets

## Updating the Deployment

### Method 1: Using Deployment Script

```bash
# From local machine
./scripts/deploy-to-ax-sas-tools.sh
```

### Method 2: Manual Update

```bash
# On server
ssh root@121.40.188.25
cd /opt/axinova-mcp-server

# Pull latest code
git pull origin main

# Rebuild image
docker build -t axinova-mcp-server:latest .

# Restart service
docker-compose down
docker-compose up -d
```

## Monitoring and Logs

### View Container Logs

```bash
# Real-time logs
docker-compose logs -f mcp-server

# Last 100 lines
docker-compose logs --tail=100 mcp-server

# Since specific time
docker-compose logs --since=1h mcp-server
```

### Check Container Status

```bash
docker-compose ps
docker-compose top
docker stats axinova-mcp-server
```

### Health Check

```bash
# Test initialization
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}' | docker-compose exec -T mcp-server /app/axinova-mcp-server

# List tools
(echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-11-25","capabilities":{},"clientInfo":{"name":"test","version":"1.0.0"}}}'; \
 echo '{"jsonrpc":"2.0","method":"initialized"}'; \
 echo '{"jsonrpc":"2.0","id":2,"method":"tools/list"}') | docker-compose exec -T mcp-server /app/axinova-mcp-server 2>/dev/null | grep -o '"name":"[^"]*"' | head -10
```

## Troubleshooting

### Container Won't Start

```bash
# Check logs
docker-compose logs mcp-server

# Common issues:
# 1. Missing .env file
# 2. Invalid API tokens
# 3. Port conflicts
```

### Connection Errors

```bash
# Test network connectivity
docker-compose exec mcp-server ping -c 3 portainer.axinova-internal.xyz
docker-compose exec mcp-server wget -O- https://prometheus.axinova-internal.xyz/api/v1/status/config
```

### TLS Certificate Errors

If you see TLS errors, ensure `config/prod.yaml` has:

```yaml
tls:
  skip_verify: true
```

### Authentication Errors

```bash
# Test tokens manually
curl -H "X-API-Key: YOUR_PORTAINER_TOKEN" https://portainer.axinova-internal.xyz/api/endpoints
curl -H "Authorization: Bearer YOUR_GRAFANA_TOKEN" https://grafana.axinova-internal.xyz/api/health
```

## Security Considerations

### API Token Management

- **Never commit tokens to git**
- Store tokens in `.env` file (gitignored)
- Rotate tokens periodically
- Use minimum required permissions

### Network Security

- MCP server runs in docker network
- No external ports exposed (stdio transport)
- All service communication over HTTPS
- TLS termination at Traefik level

### Access Control

- SSH access restricted to authorized users
- Docker daemon socket not exposed
- Container runs as non-root user (`mcp:mcp`)

## Integration with Claude Desktop

After deployment, configure Claude Desktop on your local machine:

### macOS/Linux

Edit `~/Library/Application Support/Claude/claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "axinova": {
      "command": "ssh",
      "args": [
        "root@121.40.188.25",
        "docker",
        "exec",
        "-i",
        "axinova-mcp-server",
        "/app/axinova-mcp-server"
      ],
      "env": {}
    }
  }
}
```

### Testing

In Claude Desktop, try:

```
List all Docker containers in Portainer
```

```
Show me Prometheus metrics for CPU usage
```

```
List my Vikunja projects
```

## Rollback Procedure

If deployment fails or causes issues:

```bash
# Stop current deployment
docker-compose down

# Revert to previous version
git log --oneline  # Find previous commit
git reset --hard <previous-commit>

# Rebuild and restart
docker build -t axinova-mcp-server:latest .
docker-compose up -d
```

## Backup and Restore

### Backup Configuration

```bash
# Backup .env file
cp /opt/axinova-mcp-server/.env /opt/axinova-mcp-server/.env.backup.$(date +%Y%m%d)

# Backup entire directory
tar -czf /opt/backups/axinova-mcp-server-$(date +%Y%m%d).tar.gz /opt/axinova-mcp-server
```

### Restore Configuration

```bash
# Restore .env
cp /opt/axinova-mcp-server/.env.backup.20260119 /opt/axinova-mcp-server/.env

# Restore from archive
tar -xzf /opt/backups/axinova-mcp-server-20260119.tar.gz -C /
```

## Performance Tuning

### Resource Limits

Edit `docker-compose.yml` to add resource constraints:

```yaml
services:
  mcp-server:
    # ... existing config ...
    deploy:
      resources:
        limits:
          cpus: '0.5'
          memory: 512M
        reservations:
          cpus: '0.1'
          memory: 128M
```

### Logging

Configure logging driver:

```yaml
services:
  mcp-server:
    # ... existing config ...
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## Next Steps

1. **Configure Tokens**: See [scripts/get-tokens.md](scripts/get-tokens.md)
2. **Test Tools**: See [TESTING.md](TESTING.md)
3. **Integrate with Claude**: Configure Claude Desktop
4. **Monitor**: Set up log monitoring and alerts
5. **Automate**: Add to CI/CD pipeline

## Support

For issues:
- GitHub: https://github.com/axinova-ai/axinova-mcp-server-go/issues
- Internal: Axinova DevOps team
- Logs: `/opt/axinova-mcp-server/` on server

## Changelog

- **2026-01-19**: Initial deployment to ax-sas-tools
  - Go 1.23 with Chinese proxy support
  - TLS skip verification for internal services
  - 40+ tools across 5 services
