# Final Setup Steps - Portainer and ax-sas-tools

**Date:** 2026-01-20
**Status:** Port 9001 opened ✅ | Manual registration needed ⏳

---

## ✅ Confirmed Working

### Port 9001 is Now Open

All 4 VPC machines are now accessible on port 9001:
- ✅ ax-dev-app (172.18.80.46:9001)
- ✅ ax-dev-db (172.18.80.47:9001)
- ✅ ax-prod-app (172.18.80.48:9001)
- ✅ ax-prod-db (172.18.80.49:9001)

**Verified:** Portainer agents are running with TLS enabled (self-signed certificates)

---

## Step 1: Register Portainer Agents in UI

The Portainer API in version 2.20.3 has validation issues that prevent programmatic registration. Manual registration via UI is required.

### Instructions

1. **Access Portainer**
   - URL: https://portainer.axinova-internal.xyz
   - Login: admin / 123321

2. **Add Each Environment**
   - Click **"Environments"** in the left sidebar
   - Click **"+ Add environment"** button
   - Select **"Docker Standalone"** (not Docker Swarm)
   - Select **"Agent"** as the connection method

3. **Environment Details**

For each machine, use these exact settings:

#### ax-dev-app
```
Name: ax-dev-app
Environment address: 172.18.80.46:9001
TLS: ✓ Enabled (check this box)
Skip TLS verification: ✓ Enabled (check this box)
```

#### ax-dev-db
```
Name: ax-dev-db
Environment address: 172.18.80.47:9001
TLS: ✓ Enabled
Skip TLS verification: ✓ Enabled
```

#### ax-prod-app
```
Name: ax-prod-app
Environment address: 172.18.80.48:9001
TLS: ✓ Enabled
Skip TLS verification: ✓ Enabled
```

#### ax-prod-db
```
Name: ax-prod-db
Environment address: 172.18.80.49:9001
TLS: ✓ Enabled
Skip TLS verification: ✓ Enabled
```

4. **Click "Add environment"** for each

5. **Verify Connection**
   - Each environment should show a green dot (Connected)
   - Click on any environment to view its containers
   - You should see containers like: portainer_agent, promtail, node-exporter

---

## Step 2: Fix ax-sas-tools Connectivity (Optional)

ax-sas-tools (121.40.188.25) is on **Aliyun Simple Application Service**, not in the VPC (vpc-bp1p00qic1fnpx7ndahon).

### Current Status
- ❌ Port 9100 (Node Exporter): BLOCKED
- ❌ Port 9001 (Portainer Agent): BLOCKED

### Why It's Blocked

ax-sas-tools is on a different network than the other 5 ECS instances. It needs firewall rules to allow:
- **From:** ax-tools public IP (120.26.32.121)
- **To:** Ports 9100 and 9001 on ax-sas-tools (121.40.188.25)

### How to Fix (Manual - Aliyun Console)

Since ax-sas-tools is on Simple Application Service, you need to configure its firewall in the Aliyun console:

1. **Login to Aliyun Console**
   - Navigate to **Simple Application Service** (轻量应用服务器)
   - Select region: **Hangzhou** (cn-hangzhou)

2. **Find ax-sas-tools Instance**
   - Look for instance with IP: 121.40.188.25

3. **Configure Firewall Rules**
   - Click on the instance
   - Go to **"Firewall"** or **"Security"** tab
   - Add two new rules:

   **Rule 1: Node Exporter**
   ```
   Application Type: Custom TCP
   Port Range: 9100
   Protocol: TCP
   Source: 120.26.32.121/32
   Policy: Allow
   Description: Node Exporter access from ax-tools
   Priority: 100
   ```

   **Rule 2: Portainer Agent**
   ```
   Application Type: Custom TCP
   Port Range: 9001
   Protocol: TCP
   Source: 120.26.32.121/32
   Policy: Allow
   Description: Portainer agent access from ax-tools
   Priority: 100
   ```

4. **Save and Apply Rules**

### Verification After Fix

```bash
# From ax-tools, test port 9100
ssh ax-tools "curl -s http://121.40.188.25:9100/metrics | head -5"
# Should show: # HELP go_gc_duration_seconds...

# Test port 9001
ssh ax-tools "curl -k -s https://121.40.188.25:9001/ | head -5"
# Should return empty or Portainer response
```

### After Fixing ax-sas-tools Firewall

1. **Add to Portainer** (same as other machines):
   ```
   Name: ax-sas-tools
   Environment address: 121.40.188.25:9001
   TLS: ✓ Enabled
   Skip TLS verification: ✓ Enabled
   ```

2. **Prometheus will auto-detect** the new metrics source within 15 seconds

---

## Step 3: Verify Everything is Working

### Check Portainer Environments

1. Go to https://portainer.axinova-internal.xyz
2. Click "Environments"
3. You should see:
   - ✅ primary (local - ax-tools)
   - ✅ ax-dev-app (green dot)
   - ✅ ax-dev-db (green dot)
   - ✅ ax-prod-app (green dot)
   - ✅ ax-prod-db (green dot)
   - (Optional) ✅ ax-sas-tools (green dot)

### Check Prometheus Targets

```bash
ssh ax-tools 'curl -s http://172.18.80.50:9090/api/v1/targets | jq -r ".data.activeTargets[] | select(.labels.job==\"node\") | \"\(.labels.host): \(.health)\""'
```

Expected output (after fixing ax-sas-tools):
```
tools: up
dev-app: up
dev-db: up
prod-app: up
prod-db: up
sas-tools: up  # After firewall fix
```

### Check Grafana Dashboards

1. **Centralized Logging**
   - URL: https://grafana.axinova-internal.xyz/d/cfaoebjy2lszka/centralized-logging
   - Should show logs from all 6 machines
   - Test filters: select different hosts, environments, containers

2. **Node Exporter Full** (Host Metrics)
   - URL: https://grafana.axinova-internal.xyz/d/rYdddlPWk/node-exporter-full
   - Select different hosts from dropdown
   - Should show CPU, memory, disk, network metrics
   - Currently shows 5 hosts (6 after fixing ax-sas-tools)

---

## Troubleshooting

### Portainer Agent Shows "Disconnected"

1. **Check agent is running:**
   ```bash
   ssh <host> "docker ps | grep portainer_agent"
   ```

2. **Check logs:**
   ```bash
   ssh <host> "docker logs portainer_agent --tail 50"
   ```

3. **Restart agent:**
   ```bash
   ssh <host> "docker restart portainer_agent"
   ```

4. **In Portainer UI:**
   - Go to Environments
   - Click on the disconnected environment
   - Click "Update environment"
   - Verify address is correct: `172.18.80.XX:9001`
   - Ensure TLS is enabled
   - Ensure "Skip TLS verification" is enabled
   - Click "Update environment"

### Node Exporter Shows "Down" in Prometheus

1. **Check if Node Exporter is running:**
   ```bash
   ssh <host> "docker ps | grep node-exporter"
   ```

2. **Test metrics endpoint:**
   ```bash
   ssh <host> "curl -s http://localhost:9100/metrics | head -5"
   ```

3. **Check firewall from ax-tools:**
   ```bash
   ssh ax-tools "curl -s http://<internal-ip>:9100/metrics | head -5"
   ```

4. **Restart Prometheus to reload config:**
   ```bash
   ssh ax-tools "docker restart observability_prometheus_1"
   ```

---

## Summary of Network Architecture

### VPC Network (vpc-bp1p00qic1fnpx7ndahon)

**Machines in VPC:** 5 ECS instances
- ax-tools (172.18.80.50) - Hub
- ax-dev-app (172.18.80.46)
- ax-dev-db (172.18.80.47)
- ax-prod-app (172.18.80.48)
- ax-prod-db (172.18.80.49)

**Communication:** Direct via internal IPs (172.18.80.x)

### External Network

**Machine:** 1 Simple Application Service instance
- ax-sas-tools (121.40.188.25)

**Communication:** Via public IP only
- From ax-tools public IP (120.26.32.121)
- To ax-sas-tools public IP (121.40.188.25)
- Requires specific firewall rules for ports 9100 and 9001

---

## Current Infrastructure Status

| Component | Status | Notes |
|-----------|--------|-------|
| Portainer Agents Deployed | ✅ 5/5 | All running with TLS |
| Port 9001 Access (VPC) | ✅ Open | Security groups configured |
| Port 9001 Access (ax-sas-tools) | ❌ Blocked | Needs SAS firewall config |
| Node Exporter Deployed | ✅ 6/6 | All running |
| Port 9100 Access (VPC) | ✅ Open | All 5 machines reachable |
| Port 9100 Access (ax-sas-tools) | ❌ Blocked | Needs SAS firewall config |
| Prometheus Targets | ✅ 5/6 up | ax-sas-tools blocked |
| Loki Log Collection | ✅ 6/6 | All shipping logs |
| Grafana Dashboards | ✅ Working | Logging + Node Exporter |
| Duplicate Stack Cleanup | ✅ Done | ax-sas-tools cleaned |

---

## Next Actions (In Order)

1. ✅ **DONE:** Port 9001 opened in security groups
2. ⏳ **NOW:** Register 4 Portainer agents in UI (5-10 minutes)
   - Use internal IPs: 172.18.80.46/47/48/49
   - Enable TLS with skip verification
3. 🔧 **Optional:** Fix ax-sas-tools firewall (5 minutes)
   - Opens ports 9100 and 9001
   - Adds 6th machine to monitoring
   - Register 5th Portainer agent

---

## Final Result

After completing all steps, you will have:

✅ **Centralized Observability**
- Single Grafana instance with all data
- Single Prometheus instance scraping all hosts
- Single Loki instance collecting all logs

✅ **Container Management**
- Portainer managing containers across all machines
- Single UI to view/manage all Docker containers

✅ **Monitoring**
- Host metrics from 6 machines (after fixing ax-sas-tools)
- Container logs from 6 machines
- Interactive dashboards with filters

✅ **Clean Architecture**
- No duplicate services
- Agents only on managed machines
- Centralized hub on ax-tools

---

## Files Reference

- `FINAL_SETUP_STEPS.md` - This file
- `SECURITY_GROUP_STATUS.md` - Security group details
- `PORTAINER_AGENT_SETUP.md` - Agent deployment info
- `DEPLOYMENT_COMPLETE_SUMMARY.md` - Full deployment summary

---

## Contact & Support

**Grafana:** https://grafana.axinova-internal.xyz (admin:123321)
**Portainer:** https://portainer.axinova-internal.xyz (admin:123321)
**Prometheus:** https://prometheus.axinova-internal.xyz

All documentation saved in: `/Users/weixia/axinova/axinova-mcp-server-go/`
