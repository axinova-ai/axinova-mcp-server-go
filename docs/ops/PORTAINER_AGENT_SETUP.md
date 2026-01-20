# Portainer Agent Setup Guide

**Date:** 2026-01-20
**Status:** ✅ Agents Deployed, ❌ Port 9001 Blocked

---

## ⚠️ IMPORTANT: Port 9001 Must Be Opened First

Before registering agents in Portainer UI, you must configure Aliyun security groups to allow port 9001.

**Current Issue:** Port 9001 is blocked by firewall, causing timeout errors:
```
Get "https://120.26.30.40:9001/ping": context deadline exceeded
```

**Solution:** See `SECURITY_GROUP_STATUS.md` for detailed instructions on opening port 9001.

---

## Completed Steps

All 5 Portainer agents have been successfully deployed and are running:

| Machine | IP | Port | Container ID | Status |
|---------|------------|------|--------------|--------|
| ax-dev-app | 120.26.30.40 | 9001 | caf0ff347935 | ✅ Running |
| ax-dev-db | 172.18.80.47 | 9001 | 2aa2609c8ef1 | ✅ Running |
| ax-prod-app | 114.55.132.190 | 9001 | d1381706b443 | ✅ Running |
| ax-prod-db | 172.18.80.49 | 9001 | 1fb533b18cf9 | ✅ Running |
| ax-sas-tools | 121.40.188.25 | 9001 | 8d3344b690d9 | ✅ Running |

All agents are using the image from the private registry: `registry.axinova-internal.xyz/portainer/agent:latest`

---

## Next Step: Register Endpoints in Portainer UI

The agents are deployed but need to be registered in the Portainer UI at https://portainer.axinova-internal.xyz

### Manual Registration Steps

1. **Access Portainer UI**
   - URL: https://portainer.axinova-internal.xyz
   - Login with admin credentials (admin:123321)

2. **Add Each Environment**
   - Click "Environments" in the left sidebar
   - Click "+ Add environment"
   - Select "Agent" as the environment type
   - **IMPORTANT:** Uncheck "TLS" (agents are running HTTP, not HTTPS)
   - **Use internal IPs** from ax-tools (172.18.80.x for VPC machines)
   - Fill in the details:

#### ax-dev-app
```
Name: ax-dev-app
Environment URL: 172.18.80.46:9001
TLS: OFF (uncheck)
```

#### ax-dev-db
```
Name: ax-dev-db
Environment URL: 172.18.80.47:9001
TLS: OFF (uncheck)
```

#### ax-prod-app
```
Name: ax-prod-app
Environment URL: 172.18.80.48:9001
TLS: OFF (uncheck)
```

#### ax-prod-db
```
Name: ax-prod-db
Environment URL: 172.18.80.49:9001
TLS: OFF (uncheck)
```

#### ax-sas-tools
```
Name: ax-sas-tools
Environment URL: 121.40.188.25:9001
TLS: OFF (uncheck)
```

3. **Verify Connection**
   - After adding each environment, verify it shows as "Connected" (green dot)
   - Click on each environment to view containers

---

## Agent Configuration

All agents were deployed with:
- **Port**: 9001 (standard Portainer agent port)
- **Restart Policy**: always
- **Socket Mount**: `/var/run/docker.sock`
- **Volumes Mount**: `/var/lib/docker/volumes`

---

## Verification Commands

To verify agents are running on each machine:

```bash
# Check all agents
for host in ax-dev-app ax-dev-db ax-prod-app ax-prod-db ax-sas-tools; do
  echo "=== $host ==="
  ssh $host "docker ps | grep portainer_agent"
  echo ""
done
```

Expected output: Each should show a running container on port 9001.

---

## Troubleshooting

### Agent Not Connecting

If an agent shows as disconnected in the UI:

1. **Check if agent is running:**
   ```bash
   ssh <host> "docker ps | grep portainer_agent"
   ```

2. **Check agent logs:**
   ```bash
   ssh <host> "docker logs portainer_agent"
   ```

3. **Restart agent:**
   ```bash
   ssh <host> "docker restart portainer_agent"
   ```

### Port Accessibility

Verify port 9001 is accessible from ax-tools (where Portainer server runs):

```bash
ssh ax-tools "telnet 120.26.30.40 9001"
```

### Firewall Rules

Ensure security groups allow TCP port 9001 from ax-tools (120.26.32.121) to all agent hosts.

---

## API Registration (Alternative Method)

**Note:** Attempted to register via API but encountered validation errors with Portainer 2.20.3 API.

If API registration is needed, the correct endpoint format may have changed. Check:
- [Portainer API Documentation](https://docs.portainer.io/api/docs)
- API version: `/api/system/status`

---

## Next Tasks

Once agents are registered in the UI:

1. ✅ Agents deployed and running
2. ⏳ **Manual registration in Portainer UI** (current step)
3. ⏳ Deploy Promtail to all 6 machines
4. ⏳ Create centralized logging dashboard in Grafana
5. ⏳ Create host monitoring dashboards

---

## Contact

- **Portainer URL**: https://portainer.axinova-internal.xyz
- **Admin Credentials**: admin:123321
- **API Token**: ptr_ChiXtsrSJZPSHRE1LAdSiBPobYttxre+ydGYimMYNyA=
