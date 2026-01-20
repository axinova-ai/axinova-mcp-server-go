# Configuration Validation Guide

This document validates the proposed configuration approach against the actual implementation.

## ❌ Incorrect Approach (Common Mistakes)

```bash
# DON'T USE THESE VARIABLE NAMES:
GITHUB_TOKEN=ghp_xxx                    # ❌ No GitHub integration
GRAFANA_API_TOKEN=glsa_xxx              # ❌ Wrong naming convention
PORTAINER_API_TOKEN=ptr_xxx             # ❌ Wrong naming convention
VIKUNJA_API_TOKEN=vk_xxx                # ❌ Wrong naming convention
```

## ✅ Correct Approach (Actual Implementation)

```bash
# USE THESE VARIABLE NAMES (APP_ prefix, __ for nesting):
ENV=prod                                        # ✅ Set environment

# Portainer
APP_PORTAINER__URL=https://portainer.axinova-internal.xyz
APP_PORTAINER__TOKEN=ptr_xxx                    # ✅ Correct naming

# Grafana
APP_GRAFANA__URL=https://grafana.axinova-internal.xyz
APP_GRAFANA__TOKEN=glsa_xxx                     # ✅ Correct naming

# Prometheus (no token needed)
APP_PROMETHEUS__URL=https://prometheus.axinova-internal.xyz

# SilverBullet
APP_SILVERBULLET__URL=https://silverbullet.axinova-internal.xyz
APP_SILVERBULLET__TOKEN=                        # ✅ Optional (leave empty if no auth)

# Vikunja
APP_VIKUNJA__URL=https://vikunja.axinova-internal.xyz
APP_VIKUNJA__TOKEN=tk_xxx                       # ✅ Correct naming
```

## Why APP_ Prefix and Double Underscore?

This server uses **Koanf** for configuration management, which requires:

1. **APP_ prefix**: All environment variables must start with `APP_`
2. **Double underscore (__)**: Represents nested structure
   - `APP_PORTAINER__TOKEN` → `portainer.token` in config
   - `APP_GRAFANA__URL` → `grafana.url` in config

## Configuration Precedence

1. `config/base.yaml` - Base defaults
2. `config/prod.yaml` - Production overrides
3. `.env` file - Environment-specific secrets (highest priority)

Example:
```yaml
# config/prod.yaml already sets:
tls:
  skip_verify: true

# You only need to set in .env:
APP_PORTAINER__TOKEN=your-token
```

## Services Integrated

| Service | Purpose | Auth Required | Variable |
|---------|---------|---------------|----------|
| **Portainer** | Container management | ✅ Yes | `APP_PORTAINER__TOKEN` |
| **Grafana** | Dashboards/metrics | ✅ Yes | `APP_GRAFANA__TOKEN` |
| **Prometheus** | Metrics queries | ❌ No | (none) |
| **SilverBullet** | Notes/wiki | ⚠️ Optional | `APP_SILVERBULLET__TOKEN` |
| **Vikunja** | Task management | ✅ Yes | `APP_VIKUNJA__TOKEN` |

**Note:** No GitHub integration in this implementation.

## Step-by-Step Configuration

### Method 1: Interactive Script (Recommended)

```bash
# On server
ssh root@121.40.188.25
cd /opt/axinova-mcp-server
bash scripts/configure-tokens.sh --local
```

### Method 2: Manual Configuration

```bash
# On server
ssh root@121.40.188.25
cd /opt/axinova-mcp-server
nano .env
```

Then paste:
```bash
ENV=prod

APP_PORTAINER__URL=https://portainer.axinova-internal.xyz
APP_PORTAINER__TOKEN=your-portainer-token-here

APP_GRAFANA__URL=https://grafana.axinova-internal.xyz
APP_GRAFANA__TOKEN=your-grafana-token-here

APP_PROMETHEUS__URL=https://prometheus.axinova-internal.xyz

APP_SILVERBULLET__URL=https://silverbullet.axinova-internal.xyz
APP_SILVERBULLET__TOKEN=

APP_VIKUNJA__URL=https://vikunja.axinova-internal.xyz
APP_VIKUNJA__TOKEN=your-vikunja-token-here
```

Save and set permissions:
```bash
chmod 600 .env
```

## Validation Checklist

Before starting the service, verify:

- [ ] `.env` file exists in `/opt/axinova-mcp-server/`
- [ ] `ENV=prod` is set (not `dev`)
- [ ] All token variables use `APP_` prefix
- [ ] All token variables use `__` (double underscore) for nesting
- [ ] File permissions are `600` (secure)
- [ ] No GitHub tokens (not used in this implementation)
- [ ] At least Portainer, Grafana, and Vikunja tokens are configured

## Common Errors and Solutions

### Error: "missing token"
```
✗ Solution: Check variable naming uses APP_ prefix and __ (not _)
```

### Error: "401 unauthorized"
```
✗ Solution: Token is incorrect or expired. Regenerate token.
```

### Error: "config not loading environment variables"
```
✗ Solution: Ensure variables start with APP_ prefix exactly.
```

### Error: "x509: certificate signed by unknown authority"
```
✗ Solution: Verify config/prod.yaml has tls.skip_verify: true
```

## Testing Configuration

After configuring .env, test without starting docker-compose:

```bash
# Test environment variable loading
cd /opt/axinova-mcp-server
source .env
echo $ENV                    # Should print: prod
echo $APP_PORTAINER__TOKEN   # Should print your token
```

## Security Best Practices

### 1. File Permissions
```bash
chmod 600 .env              # Only owner can read/write
chown root:root .env        # Owned by root
```

### 2. Token Management
- Never commit .env to git (already in .gitignore)
- Rotate tokens periodically
- Use minimum required permissions (read-only where possible)
- Store backups securely

### 3. Backup Configuration
```bash
# Before changes
cp .env .env.backup.$(date +%Y%m%d)

# Secure backup location
mkdir -p /opt/backups/mcp-server
cp .env /opt/backups/mcp-server/.env.$(date +%Y%m%d)
chmod 600 /opt/backups/mcp-server/.env.*
```

## Next Steps After Configuration

1. **Start service:**
   ```bash
   docker-compose up -d
   ```

2. **Verify startup:**
   ```bash
   docker-compose logs -f mcp-server
   ```

3. **Test tools:**
   See TESTING.md for comprehensive testing procedures.

4. **Monitor:**
   ```bash
   docker-compose ps
   docker stats axinova-mcp-server
   ```

## References

- Configuration system: `internal/config/config.go`
- Environment example: `.env.example`
- Token generation: `scripts/get-tokens.md`
- Full deployment: `DEPLOYMENT.md`
- Testing guide: `TESTING.md`
