# How to Get API Tokens

This guide explains how to obtain API tokens for each service.

## 1. Portainer

### Method 1: Via Web UI
1. Login to https://portainer.axinova-internal.xyz
2. Click your username (top right) → **My account**
3. Go to **Access tokens** tab
4. Click **+ Add access token**
5. Name: `mcp-server`
6. Click **Add access token**
7. **Copy the token immediately** (it won't be shown again)

### Method 2: Via API
```bash
# Login first
PORTAINER_USER="admin"
PORTAINER_PASS="your-password"

# Get JWT token
JWT=$(curl -s -X POST https://portainer.axinova-internal.xyz/api/auth \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$PORTAINER_USER\",\"password\":\"$PORTAINER_PASS\"}" \
  | jq -r '.jwt')

# Create API token
curl -X POST https://portainer.axinova-internal.xyz/api/users/admin/tokens \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Type: application/json" \
  -d '{"description":"mcp-server"}' | jq
```

---

## 2. Grafana

### Method 1: Via Web UI
1. Login to https://grafana.axinova-internal.xyz
2. Go to **Configuration** (gear icon) → **API Keys**
3. Click **New API Key**
4. Key name: `mcp-server`
5. Role: **Admin** or **Editor**
6. Time to live: Leave empty or set expiration
7. Click **Add**
8. **Copy the API key immediately**

### Method 2: Via API
```bash
GRAFANA_USER="admin"
GRAFANA_PASS="your-password"

curl -X POST https://grafana.axinova-internal.xyz/api/auth/keys \
  -u "$GRAFANA_USER:$GRAFANA_PASS" \
  -H "Content-Type: application/json" \
  -d '{"name":"mcp-server","role":"Admin"}' | jq
```

### Method 3: Using Service Account (Recommended for Grafana 9+)
```bash
# Create service account
curl -X POST https://grafana.axinova-internal.xyz/api/serviceaccounts \
  -u "$GRAFANA_USER:$GRAFANA_PASS" \
  -H "Content-Type: application/json" \
  -d '{"name":"mcp-server","role":"Admin"}' | jq

# Create token for service account (use ID from above)
SA_ID=1
curl -X POST https://grafana.axinova-internal.xyz/api/serviceaccounts/$SA_ID/tokens \
  -u "$GRAFANA_USER:$GRAFANA_PASS" \
  -H "Content-Type: application/json" \
  -d '{"name":"mcp-server-token"}' | jq
```

---

## 3. Prometheus

**No authentication required** for internal access. Prometheus is typically accessed through Grafana or without authentication in internal networks.

If you have authentication enabled, you'll need to configure it separately.

---

## 4. SilverBullet

SilverBullet typically uses one of these methods:

### Method 1: Token Authentication
Check your SilverBullet configuration file for token settings:
```bash
# On ax-sas-tools
cat /opt/silverbullet/.env
# or
docker exec silverbullet cat /config/.env
```

### Method 2: No Authentication
If SilverBullet is configured without authentication (common for internal deployments), you can leave the token empty:
```bash
APP_SILVERBULLET__TOKEN=
```

### Method 3: Create Token
If using authentication, create a token in SilverBullet:
1. Check SilverBullet documentation for token generation
2. May require editing configuration file directly

---

## 5. Vikunja

### Method 1: Via Web UI
1. Login to https://vikunja.axinova-internal.xyz
2. Click your avatar (top right) → **Settings**
3. Go to **API Tokens** or **Tokens**
4. Click **Create a new token**
5. Name: `mcp-server`
6. Permissions: Full access
7. Click **Create**
8. **Copy the token**

### Method 2: Via API
```bash
VIKUNJA_USER="your-username"
VIKUNJA_PASS="your-password"

# Login to get JWT token
JWT=$(curl -s -X POST https://vikunja.axinova-internal.xyz/api/v1/login \
  -H "Content-Type: application/json" \
  -d "{\"username\":\"$VIKUNJA_USER\",\"password\":\"$VIKUNJA_PASS\"}" \
  | jq -r '.token')

# Create API token
curl -X PUT https://vikunja.axinova-internal.xyz/api/v1/tokens \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Type: application/json" \
  -d '{"title":"mcp-server"}' | jq
```

---

## Quick Setup Script

Save this as `scripts/setup-tokens.sh`:

```bash
#!/bin/bash

echo "=== MCP Server Token Setup ==="
echo ""

# Function to get token securely
get_token() {
  local service=$1
  local token
  echo -n "$service token: "
  read -s token
  echo ""
  echo "$token"
}

# Backup existing .env
if [ -f .env ]; then
  cp .env .env.backup
  echo "✓ Backed up existing .env to .env.backup"
fi

# Copy from example
cp .env.example .env

echo ""
echo "Please enter your API tokens (input hidden):"
echo ""

PORTAINER_TOKEN=$(get_token "Portainer")
GRAFANA_TOKEN=$(get_token "Grafana")
SILVERBULLET_TOKEN=$(get_token "SilverBullet")
VIKUNJA_TOKEN=$(get_token "Vikunja")

# Update .env file
sed -i "s|APP_PORTAINER__TOKEN=.*|APP_PORTAINER__TOKEN=$PORTAINER_TOKEN|" .env
sed -i "s|APP_GRAFANA__TOKEN=.*|APP_GRAFANA__TOKEN=$GRAFANA_TOKEN|" .env
sed -i "s|APP_SILVERBULLET__TOKEN=.*|APP_SILVERBULLET__TOKEN=$SILVERBULLET_TOKEN|" .env
sed -i "s|APP_VIKUNJA__TOKEN=.*|APP_VIKUNJA__TOKEN=$VIKUNJA_TOKEN|" .env

echo ""
echo "✓ Tokens configured in .env"
echo ""
echo "Next steps:"
echo "  1. Review .env file"
echo "  2. Run: make build"
echo "  3. Test: ./bin/axinova-mcp-server"
```

Make it executable:
```bash
chmod +x scripts/setup-tokens.sh
./scripts/setup-tokens.sh
```

---

## Verification

After setting tokens, verify connectivity:

```bash
# Test Portainer
curl -H "X-API-Key: YOUR_PORTAINER_TOKEN" \
  https://portainer.axinova-internal.xyz/api/endpoints

# Test Grafana
curl -H "Authorization: Bearer YOUR_GRAFANA_TOKEN" \
  https://grafana.axinova-internal.xyz/api/health

# Test Prometheus (no auth)
curl https://prometheus.axinova-internal.xyz/api/v1/status/config

# Test Vikunja
curl -H "Authorization: Bearer YOUR_VIKUNJA_TOKEN" \
  https://vikunja.axinova-internal.xyz/api/v1/user
```
