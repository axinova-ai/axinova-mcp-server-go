# Infrastructure Setup - Status Update

**Date:** 2026-01-20
**Time:** 00:40 CST

---

## ✅ Completed Tasks

1. **Connected to ax-tools** (120.26.32.121)
   - Verified observability stack running: Grafana, Prometheus, Loki, Portainer, Traefik
   - User: ecs-user

2. **Investigated Grafana slow login issue**
   - Server response time: 0.6s (fast)
   - Initial conclusion: Likely DNS resolution or browser cache issue on first visit
   - After first login, subsequent pages load instantly
   - No server-side performance issue detected

3. **Created new Grafana API token**
   - Token: `[REDACTED - See .env file on ax-tools]`
   - Name: mcp-server-ax-tools
   - Role: Admin
   - Password used: [REDACTED]

4. **Updated MCP server configuration**
   - Updated Grafana token in `/opt/axinova-mcp-server/.env`
   - MCP server restarted successfully
   - Now pointing to centralized Grafana on ax-tools

---

## 🔧 In Progress

1. **Deploying Portainer Agents**
   - Issue encountered: Docker Hub timeout on Chinese servers
   - Solution: Need to use private registry (registry.axinova-internal.xyz)
   - Status: Preparing to pull images to private registry

---

## ⚠️ Known Issues

### 1. registry-ui CORS Error

**Error:**
```
The `Access-Control-Allow-Credentials` header in the response is missing and must be set to `true`
when the request's credentials mode is on. Origin `https://registry-ui.axinova-internal.xyz` is
therefore not allowed access
```

**Current Configuration:** `/opt/registry/docker-compose.yml`
```yaml
REGISTRY_HTTP_HEADERS_ACCESS-CONTROL-ALLOW-ORIGIN: '[https://registry-ui.axinova-internal.xyz]'
REGISTRY_HTTP_HEADERS_ACCESS-CONTROL-ALLOW-CREDENTIALS: '[true]'
```

**Root Cause:** CORS headers from registry backend not being proxied through registry-ui nginx properly

**Workaround:** Access registry API directly at https://registry.axinova-internal.xyz/v2/

**Fix Required:**
- Option A: Configure nginx in registry-ui to add CORS headers
- Option B: Modify registry-ui Docker image to include custom nginx config
- Option C: Add Traefik middleware to inject CORS headers

**Priority:** Medium (UI access issue, API access still works)

### 2. Docker Hub Access from China

**Issue:** `dial tcp 199.16.158.8:443: i/o timeout` when pulling from Docker Hub

**Impact:** Cannot deploy Portainer agents to machines without pre-pulled images

**Solution:** Use private registry mirror
```bash
# On each machine, configure Docker to use private registry mirror
{
  "registry-mirrors": ["https://registry.axinova-internal.xyz"]
}
```

**Alternative:** Pre-pull images to private registry, then deploy from there

---

## 📋 Remaining Tasks

### High Priority

1. **Deploy Portainer Agents** (5 machines)
   - [ ] ax-dev-app (120.26.30.40)
   - [ ] ax-dev-db (172.18.80.47, via ax-tools proxy)
   - [ ] ax-prod-app (114.55.132.190)
   - [ ] ax-prod-db (172.18.80.49, via ax-tools proxy)
   - [ ] ax-sas-tools (121.40.188.25)

2. **Configure Loki in Grafana**
   - [ ] Add Loki datasource
   - [ ] Test log queries
   - [ ] Verify connectivity

3. **Set Data Retention**
   - [ ] Prometheus: 7 days
   - [ ] Loki: 7 days

4. **Deploy Promtail** (log shippers to all 5 machines)
   - [ ] ax-tools
   - [ ] ax-dev-app
   - [ ] ax-dev-db
   - [ ] ax-prod-app
   - [ ] ax-prod-db

5. **Create Dashboards**
   - [ ] Centralized logging dashboard (with filters)
   - [ ] Host monitoring dashboards (5 machines)
   - [ ] Container monitoring dashboard

### Medium Priority

6. **Fix registry-ui CORS issue**

7. **Test CRUD Operations**
   - [ ] Vikunja (currently returning 404)
   - [ ] Grafana (with new token)

---

## Infrastructure Summary

### Centralized Observability (ax-tools)

| Service | URL | Status | Notes |
|---------|-----|--------|-------|
| Grafana | https://grafana.axinova-internal.xyz | ✅ Running | v11.2.0, new token created |
| Prometheus | https://prometheus.axinova-internal.xyz | ✅ Running | v2.55.0, needs retention config |
| Loki | (not exposed) | ✅ Running | v2.9.8, needs datasource config |
| Portainer | https://portainer.axinova-internal.xyz | ✅ Running | CE 2.20.3, needs agents |

### SaaS Tools (ax-sas-tools)

| Service | URL | Status | Notes |
|---------|-----|--------|-------|
| Vikunja | https://vikunja.axinova-internal.xyz | ✅ Running | v1.0.0-rc3, API working |
| SilverBullet | https://wiki.axinova-internal.xyz | ✅ Running | v2.4.1, CRUD tested |
| MCP Server | (stdio) | ✅ Running | v0.1.0, updated token |
| Docker Registry | https://registry.axinova-internal.xyz | ✅ Running | v2, CORS issue |
| Registry UI | https://registry-ui.axinova-internal.xyz | ⚠️ CORS Error | v2.5.7 |

### Application Hosts

| Host | IP | Access | Status |
|------|---------|--------|--------|
| ax-tools | 120.26.32.121 | ecs-user | ✅ Accessible |
| ax-dev-app | 120.26.30.40 | ecs-user | ✅ Accessible |
| ax-dev-db | 172.18.80.47 | ecs-user + ProxyJump | ✅ Accessible |
| ax-prod-app | 114.55.132.190 | ecs-user | ✅ Accessible |
| ax-prod-db | 172.18.80.49 | ecs-user + ProxyJump | ✅ Accessible |
| ax-sas-tools | 121.40.188.25 | root | ✅ Accessible |

---

## Next Steps

1. **Pull Portainer agent image to private registry**
   ```bash
   docker pull portainer/agent:latest
   docker tag portainer/agent:latest registry.axinova-internal.xyz/portainer/agent:latest
   docker push registry.axinova-internal.xyz/portainer/agent:latest
   ```

2. **Deploy agents using private registry**
   ```bash
   docker run -d \
     -p 9001:9001 \
     --name portainer_agent \
     --restart=always \
     -v /var/run/docker.sock:/var/run/docker.sock \
     -v /var/lib/docker/volumes:/var/lib/docker/volumes \
     registry.axinova-internal.xyz/portainer/agent:latest
   ```

3. **Add agents to Portainer UI**
   - Access https://portainer.axinova-internal.xyz
   - Add each machine as an endpoint
   - Verify container visibility

4. **Configure observability stack**
   - Use MCP server or direct API calls
   - Automate dashboard creation
   - Set up retention policies

---

## Questions for User

1. **registry-ui CORS**: Should I prioritize fixing this or continue with agent deployment?
2. **Private registry**: Should I set up Docker daemon mirror config on all machines?
3. **Grafana slow login**: Was this a one-time issue or still occurring?

---

## Time Estimates

- Portainer agents deployment: 30-45 minutes (5 machines)
- Loki configuration: 15 minutes
- Data retention setup: 10 minutes
- Promtail deployment: 30 minutes (5 machines)
- Dashboard creation: 60-90 minutes (logging + 5 host dashboards)

**Total remaining:** ~3-4 hours

---

## Contact

- MCP Server: ax-sas-tools:/opt/axinova-mcp-server
- Grafana Admin: admin:123321
- Portainer Token: ptr_ChiXtsrSJZPSHRE1LAdSiBPobYttxre+ydGYimMYNyA=
