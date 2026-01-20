# Security Group Configuration Status

**Date:** 2026-01-20
**Time:** ~04:00 CST

---

## ✅ Port 9100 (Node Exporter) - FIXED

**Status:** Working for 5 out of 6 hosts

### Verification Results

```bash
$ ssh ax-tools 'curl -s http://172.18.80.50:9090/api/v1/targets | jq -r ".data.activeTargets[] | select(.labels.job==\"node\")"'

✅ tools (172.18.80.50):      UP
✅ dev-app (172.18.80.46):    UP
✅ dev-db (172.18.80.47):     UP
✅ prod-app (172.18.80.48):   UP
✅ prod-db (172.18.80.49):    UP
❌ sas-tools (121.40.188.25): DOWN
```

### Working

- All VPC-internal machines (172.18.80.x) can now be scraped by Prometheus
- Host monitoring dashboards will show metrics for 5 machines
- Security group rule successfully applied

### sas-tools (121.40.188.25) Still Down

**Reason:** Different network (not in VPC), requires public IP access

**To Fix (optional):**
Add security group rule on ax-sas-tools (121.40.188.25):
- **Protocol:** TCP
- **Port:** 9100
- **Source:** 120.26.32.121/32 (ax-tools public IP)

**Command to verify:**
```bash
ssh ax-tools "curl -s http://121.40.188.25:9100/metrics | head -5"
```

---

## ❌ Port 9001 (Portainer Agent) - BLOCKED

**Status:** All 5 agents blocked by firewall

### Error When Adding Agents in Portainer UI

```
Failure
Get "https://120.26.30.40:9001/ping": context deadline exceeded
(Client.Timeout exceeded while awaiting headers)
```

### Root Cause

Portainer server (on ax-tools, 120.26.32.121) cannot reach agents because:
1. Port 9001 is blocked by Aliyun security groups
2. Portainer is trying to use HTTPS when agents are running HTTP

### Verification

**Agents are running and listening:**
```bash
$ ssh ax-dev-app "ss -tlnp | grep 9001"
LISTEN 0      4096         0.0.0.0:9001      0.0.0.0:*
LISTEN 0      4096            [::]:9001         [::]:*
```

**Port is blocked from ax-tools:**
```bash
$ ssh ax-tools 'timeout 2 bash -c "echo > /dev/tcp/120.26.30.40/9001"'
# Timeout - port blocked
```

### Solution: Open Port 9001 in Security Groups

Add security group rules to allow Portainer communication:

#### For VPC Machines (ax-dev-app, ax-dev-db, ax-prod-app, ax-prod-db)

**Inbound Rule:**
- **Protocol:** TCP
- **Port:** 9001
- **Source:** 172.18.80.50/32 (ax-tools internal IP)
- **Description:** Portainer agent access from ax-tools

#### For ax-sas-tools (121.40.188.25)

**Inbound Rule:**
- **Protocol:** TCP
- **Port:** 9001
- **Source:** 120.26.32.121/32 (ax-tools public IP)
- **Description:** Portainer agent access from ax-tools

### Using Aliyun CLI to Verify

After configuring security groups, verify with:

```bash
# Check if port 9001 rules exist
aliyun ecs DescribeSecurityGroupAttribute \
  --SecurityGroupId <sg-id> \
  --Direction ingress \
  --NicType intranet \
  | jq '.Permissions.Permission[] | select(.PortRange == "9001/9001")'
```

### Expected Result After Fix

```bash
$ ssh ax-tools 'for ip in 172.18.80.46 172.18.80.47 172.18.80.48 172.18.80.49; do
  curl -s http://$ip:9001/ping
done'

# Should return: "Portainer agent" for each IP
```

---

## ✅ Cleanup Completed: Duplicate Observability Stack Removed

### Issue Identified

Duplicate Grafana, Prometheus, and Loki containers were running on ax-sas-tools, conflicting with the centralized approach.

### Containers Removed

```bash
✓ observability-grafana-1 (grafana/grafana:11.2.0)
✓ observability-prometheus-1 (prom/prometheus:v2.55.0)
✓ observability-loki-1 (grafana/loki:2.9.8)
```

### Remaining Containers on ax-sas-tools (Correct)

```
✓ node-exporter - Host metrics collection
✓ promtail - Log shipping to centralized Loki
✓ portainer_agent - Container management
```

### Centralized Architecture (Confirmed)

**Observability Hub (ax-tools: 120.26.32.121)**
- Grafana (UI): https://grafana.axinova-internal.xyz
- Prometheus (metrics DB): https://prometheus.axinova-internal.xyz
- Loki (logs DB): http://172.18.80.50:3100
- Portainer (container mgmt): https://portainer.axinova-internal.xyz

**All Other Machines (Agents Only)**
- Node Exporter (port 9100): Metrics collection
- Promtail (port 9080): Log shipping
- Portainer Agent (port 9001): Container management

---

## Summary of Changes

| Item | Before | After | Status |
|------|--------|-------|--------|
| Port 9100 access | Blocked | Open (VPC) | ✅ Fixed |
| Prometheus targets | 1/6 up | 5/6 up | ✅ Working |
| Port 9001 access | Blocked | Blocked | ❌ Needs fix |
| Portainer agents | Not registered | Not registered | ⏳ Blocked by port 9001 |
| Observability stack on ax-sas-tools | Duplicate running | Removed | ✅ Cleaned up |
| Architecture | Mixed | Centralized | ✅ Correct |

---

## Next Steps

### 1. Configure Security Groups for Port 9001 (HIGH PRIORITY)

**Option A: Using Aliyun Console**

1. Log in to Aliyun console
2. Navigate to ECS → Security Groups
3. For each machine (ax-dev-app, ax-dev-db, ax-prod-app, ax-prod-db):
   - Select security group
   - Add Inbound Rule:
     - Protocol: TCP
     - Port: 9001/9001
     - Source: 172.18.80.50/32 (ax-tools internal)
     - Priority: 100
     - Policy: Allow
4. For ax-sas-tools (121.40.188.25):
   - Add Inbound Rule:
     - Protocol: TCP
     - Port: 9001/9001
     - Source: 120.26.32.121/32 (ax-tools public)

**Option B: Using Aliyun CLI**

```bash
# For each VPC machine
aliyun ecs AuthorizeSecurityGroup \
  --SecurityGroupId <sg-id> \
  --IpProtocol tcp \
  --PortRange 9001/9001 \
  --SourceCidrIp 172.18.80.50/32 \
  --Priority 100 \
  --Description "Portainer agent access from ax-tools"

# For ax-sas-tools
aliyun ecs AuthorizeSecurityGroup \
  --SecurityGroupId <sg-id-sas-tools> \
  --IpProtocol tcp \
  --PortRange 9001/9001 \
  --SourceCidrIp 120.26.32.121/32 \
  --Priority 100 \
  --Description "Portainer agent access from ax-tools"
```

### 2. Verify Port 9001 Access

```bash
# From ax-tools
ssh ax-tools 'for ip in 172.18.80.46 172.18.80.47 172.18.80.48 172.18.80.49 121.40.188.25; do
  echo -n "$ip:9001 - "
  timeout 2 curl -s http://$ip:9001/ping && echo "OK" || echo "BLOCKED"
done'
```

Expected output:
```
172.18.80.46:9001 - Portainer agent OK
172.18.80.47:9001 - Portainer agent OK
172.18.80.48:9001 - Portainer agent OK
172.18.80.49:9001 - Portainer agent OK
121.40.188.25:9001 - Portainer agent OK
```

### 3. Register Portainer Agents in UI

Once port 9001 is open:

1. Access https://portainer.axinova-internal.xyz
2. Navigate to Environments → Add environment
3. Select "Agent" type
4. **Important:** Use **HTTP** (not HTTPS) and **internal IPs**:
   - ax-dev-app: `172.18.80.46:9001`
   - ax-dev-db: `172.18.80.47:9001`
   - ax-prod-app: `172.18.80.48:9001`
   - ax-prod-db: `172.18.80.49:9001`
   - ax-sas-tools: `121.40.188.25:9001`
5. Click "Add environment" for each

### 4. (Optional) Enable ax-sas-tools Metrics

If you want metrics from ax-sas-tools as well:

```bash
# Add security group rule on ax-sas-tools
aliyun ecs AuthorizeSecurityGroup \
  --SecurityGroupId <sg-id-sas-tools> \
  --IpProtocol tcp \
  --PortRange 9100/9100 \
  --SourceCidrIp 120.26.32.121/32 \
  --Priority 100 \
  --Description "Node Exporter access from ax-tools"
```

---

## Verification Commands

### Check All Services Status

```bash
# Node Exporter (metrics)
for host in ax-tools ax-dev-app ax-dev-db ax-prod-app ax-prod-db ax-sas-tools; do
  echo "=== $host ==="
  ssh $host "docker ps --filter name=node-exporter --format 'Status: {{.Status}}'"
done

# Promtail (logs)
for host in ax-tools ax-dev-app ax-dev-db ax-prod-app ax-prod-db ax-sas-tools; do
  echo "=== $host ==="
  ssh $host "docker ps --filter name=promtail --format 'Status: {{.Status}}'"
done

# Portainer Agent (container mgmt)
for host in ax-dev-app ax-dev-db ax-prod-app ax-prod-db ax-sas-tools; do
  echo "=== $host ==="
  ssh $host "docker ps --filter name=portainer_agent --format 'Status: {{.Status}}'"
done
```

### Check Prometheus Targets

```bash
ssh ax-tools 'curl -s http://172.18.80.50:9090/api/v1/targets | jq -r ".data.activeTargets[] | select(.labels.job==\"node\") | \"\(.labels.host): \(.health)\""'
```

### Check Loki Log Sources

```bash
ssh ax-tools 'curl -s "http://172.18.80.50:3100/loki/api/v1/label/host/values" | jq -r ".data[]"'
```

---

## Current Infrastructure Status

### ✅ Working Now

- **Host Monitoring:** 5/6 machines (all VPC machines)
- **Centralized Logging:** 6/6 machines (all shipping logs)
- **Grafana Dashboards:**
  - Centralized Logging: https://grafana.axinova-internal.xyz/d/cfaoebjy2lszka/centralized-logging
  - Node Exporter Full: https://grafana.axinova-internal.xyz/d/rYdddlPWk/node-exporter-full

### ⏳ Pending

- **Portainer Agent Registration:** Blocked by port 9001 (security group needs configuration)
- **ax-sas-tools Metrics:** Optional, requires security group update

---

## Files Updated

- `SECURITY_GROUP_STATUS.md` - This file
- `DEPLOYMENT_COMPLETE_SUMMARY.md` - Updated with cleanup status
- `PORTAINER_AGENT_SETUP.md` - Updated with HTTP configuration

---

## Contact

For security group configuration:
- Aliyun Console: https://ecs.console.aliyun.com/
- Aliyun CLI docs: https://www.alibabacloud.com/help/en/ecs/developer-reference/api-ecs-2014-05-26-authorizesecuritygroup
