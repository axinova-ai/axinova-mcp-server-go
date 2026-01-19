# Token Generation Walkthrough

Complete step-by-step guide to generate API tokens for all services.

## Overview

You need tokens for:
- ✅ **Portainer** (Required) - Container management
- ✅ **Grafana** (Required) - Dashboards and metrics
- ❌ **Prometheus** (Not required) - Uses Grafana's token
- ⚠️ **SilverBullet** (Optional) - May not require auth
- ✅ **Vikunja** (Required) - Task management

---

## 1. Portainer Token (REQUIRED)

### Method A: Web UI (Recommended)

1. **Open Portainer:**
   ```
   https://portainer.axinova-internal.xyz
   ```

2. **Login** with your credentials

3. **Navigate to Access Tokens:**
   - Click your username in the top-right corner
   - Select "My account"
   - Click on "Access tokens" tab

4. **Create New Token:**
   - Click "+ Add access token"
   - Description: `mcp-server`
   - Click "Add access token"

5. **Copy the Token:**
   - Token starts with `ptr_`
   - **IMPORTANT:** Copy immediately - it won't be shown again!
   - Save it temporarily in a secure note

### Method B: API (Advanced)

```bash
# 1. Login to get JWT
JWT=$(curl -s -X POST https://portainer.axinova-internal.xyz/api/auth \
  -H "Content-Type: application/json" \
  -d '{"username":"YOUR_USERNAME","password":"YOUR_PASSWORD"}' \
  | jq -r '.jwt')

# 2. Get your user ID
USER_ID=$(curl -s -H "Authorization: Bearer $JWT" \
  https://portainer.axinova-internal.xyz/api/users/me \
  | jq -r '.Id')

# 3. Create API token
curl -X POST https://portainer.axinova-internal.xyz/api/users/$USER_ID/tokens \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Type: application/json" \
  -d '{"description":"mcp-server"}' \
  | jq -r '.rawAPIKey'
```

**Result:** Token like `ptr_xxxxxxxxxxxxxxxxxxxxxxxxxx`

---

## 2. Grafana Token (REQUIRED)

### Method A: Web UI (Recommended)

1. **Open Grafana:**
   ```
   https://grafana.axinova-internal.xyz
   ```

2. **Login** with your credentials

3. **Navigate to API Keys:**
   - Click the gear icon (⚙️) in the left sidebar
   - Select "API Keys" (under Configuration)

4. **Create New API Key:**
   - Click "+ New API Key" or "Add API key"
   - **Key name:** `mcp-server`
   - **Role:** Admin (required for full MCP functionality)
   - **Time to live:** Leave empty (no expiration) or set as needed

5. **Copy the Token:**
   - Token starts with `glsa_` or `eyJ`
   - **IMPORTANT:** Copy immediately - it won't be shown again!
   - Save it temporarily

### Method B: API (Advanced)

```bash
# Using admin credentials
GRAFANA_USER="admin"
GRAFANA_PASS="your-password"

curl -X POST https://grafana.axinova-internal.xyz/api/auth/keys \
  -u "$GRAFANA_USER:$GRAFANA_PASS" \
  -H "Content-Type: application/json" \
  -d '{"name":"mcp-server","role":"Admin"}' \
  | jq -r '.key'
```

### Method C: Service Account (Grafana 9+, Recommended)

```bash
# 1. Create service account
SA_ID=$(curl -X POST https://grafana.axinova-internal.xyz/api/serviceaccounts \
  -u "$GRAFANA_USER:$GRAFANA_PASS" \
  -H "Content-Type: application/json" \
  -d '{"name":"mcp-server","role":"Admin"}' \
  | jq -r '.id')

# 2. Create token for service account
curl -X POST https://grafana.axinova-internal.xyz/api/serviceaccounts/$SA_ID/tokens \
  -u "$GRAFANA_USER:$GRAFANA_PASS" \
  -H "Content-Type: application/json" \
  -d '{"name":"mcp-server-token"}' \
  | jq -r '.key'
```

**Result:** Token like `glsa_xxxxxxxxxxxxxxxxxxxxxxxxxxxx` or `eyJrIjoixxxxxx...`

---

## 3. Prometheus (NO TOKEN NEEDED)

Prometheus is typically accessed without authentication in internal networks.

The MCP server will access Prometheus through Grafana's datasource query API, which uses the Grafana token.

**No action required for Prometheus.**

---

## 4. SilverBullet Token (OPTIONAL)

SilverBullet may or may not require authentication depending on your setup.

### Check if Authentication is Required:

```bash
# Test access without auth
curl -I https://silverbullet.axinova-internal.xyz/index.json
```

- **200 OK** → No auth required, leave token empty
- **401 Unauthorized** → Auth required, follow below

### If Authentication is Required:

**SilverBullet typically uses token-based auth configured in its settings.**

1. **Check SilverBullet Configuration:**
   ```bash
   # On ax-sas-tools server
   ssh root@121.40.188.25
   docker exec silverbullet cat /config/.env
   # or
   cat /path/to/silverbullet/config
   ```

2. **Look for token configuration:**
   - Variable might be `SB_AUTH_TOKEN` or similar
   - Token might be in space configuration

3. **If no token exists:**
   - Leave `APP_SILVERBULLET__TOKEN` empty in .env
   - The MCP server will work if SilverBullet has no auth

**For now, you can skip this and leave empty.**

---

## 5. Vikunja Token (REQUIRED)

### Method A: Web UI (Recommended)

1. **Open Vikunja:**
   ```
   https://vikunja.axinova-internal.xyz
   ```

2. **Login** with your credentials

3. **Navigate to Settings:**
   - Click your avatar/name in the top-right
   - Select "Settings"

4. **Go to API Tokens:**
   - Look for "API Tokens" or "Tokens" in the left menu
   - If not visible, check under "Advanced" or "Developer"

5. **Create New Token:**
   - Click "Create a new token" or "+ New Token"
   - **Title/Name:** `mcp-server`
   - **Permissions:** Full access / All permissions
   - Click "Create" or "Generate"

6. **Copy the Token:**
   - Token format varies (might start with `tk_` or be a JWT)
   - **IMPORTANT:** Copy immediately!
   - Save it temporarily

### Method B: API (Advanced)

```bash
# 1. Login to get JWT
JWT=$(curl -s -X POST https://vikunja.axinova-internal.xyz/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"username":"YOUR_USERNAME","password":"YOUR_PASSWORD"}' \
  | jq -r '.token')

# 2. Create API token
curl -X PUT https://vikunja.axinova-internal.xyz/api/v1/tokens \
  -H "Authorization: Bearer $JWT" \
  -H "Content-Type: application/json" \
  -d '{"title":"mcp-server"}' \
  | jq -r '.token'
```

**Result:** Token (format varies by Vikunja version)

---

## Token Verification

After collecting tokens, verify they work:

### Test Portainer:
```bash
curl -H "X-API-Key: YOUR_PORTAINER_TOKEN" \
  https://portainer.axinova-internal.xyz/api/endpoints
```
Expected: JSON list of endpoints (not 401 error)

### Test Grafana:
```bash
curl -H "Authorization: Bearer YOUR_GRAFANA_TOKEN" \
  https://grafana.axinova-internal.xyz/api/health
```
Expected: `{"commit":"...","database":"ok",..."}`

### Test Vikunja:
```bash
curl -H "Authorization: Bearer YOUR_VIKUNJA_TOKEN" \
  https://vikunja.axinova-internal.xyz/api/v1/user
```
Expected: JSON with user info (not 401 error)

---

## Configure .env with Tokens

Once you have the tokens, update `.env`:

```bash
ssh root@121.40.188.25
cd /opt/axinova-mcp-server
nano .env
```

Update these lines:
```bash
APP_PORTAINER__TOKEN=ptr_your_actual_token_here
APP_GRAFANA__TOKEN=glsa_your_actual_token_here
APP_SILVERBULLET__TOKEN=                           # Leave empty if no auth
APP_VIKUNJA__TOKEN=your_actual_token_here
```

Save and set secure permissions:
```bash
chmod 600 .env
chown root:root .env
```

---

## Quick Reference

| Service | Token Starts With | Where to Get |
|---------|-------------------|--------------|
| Portainer | `ptr_` | User Settings → Access Tokens |
| Grafana | `glsa_` or `eyJ` | Configuration → API Keys |
| Prometheus | N/A | Not needed |
| SilverBullet | varies | Optional - check config |
| Vikunja | varies | Settings → API Tokens |

---

## Troubleshooting

### "I don't see API Keys in Grafana"
- Your user might not have admin permissions
- Try: Configuration (gear) → Service Accounts → Create Service Account

### "Portainer shows 'Access Denied'"
- You need admin or at least environment admin role
- Contact Portainer administrator to grant permissions

### "Vikunja doesn't have API Tokens section"
- Older Vikunja versions might not support API tokens
- Alternative: Use JWT token from login (Method B)

### "Token expired or invalid"
- Regenerate the token
- Check for copy/paste errors (no extra spaces)
- Ensure token has correct permissions

---

## Security Notes

- **Never share tokens** - they grant full access to services
- **Store securely** - use password manager temporarily
- **Rotate regularly** - create new tokens periodically
- **Minimum permissions** - use read-only where possible (after testing)
- **Monitor usage** - check service logs for suspicious activity

---

## Next Steps

After collecting all tokens:

1. Configure .env (see above)
2. Start the service: `docker-compose up -d`
3. Check logs: `docker-compose logs -f mcp-server`
4. Test tools: See TESTING.md

For help, see:
- DEPLOYMENT.md
- VALIDATION.md
- scripts/get-tokens.md
