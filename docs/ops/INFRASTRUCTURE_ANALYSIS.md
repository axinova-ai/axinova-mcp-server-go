# Axinova Infrastructure Analysis & Recommendations

**Date:** 2026-01-19
**Analyst:** MCP Server Deployment Team

---

## Executive Summary

This document analyzes the current infrastructure setup, identifies issues, and provides actionable recommendations for:
1. Observability architecture consolidation
2. Portainer agent deployment
3. Centralized logging with Loki
4. Monitoring dashboards with Grafana

---

## Current Infrastructure Inventory

### Machines

| Machine | IP | Role | Observations |
|---------|---------|------|--------------|
| ax-tools | 121.40.188.18 | Centralized Observability (expected) | SSH connection failed - needs investigation |
| ax-sas-tools | 121.40.188.25 | SaaS Tools + Local Observability | Has full observability stack running |
| ax-dev-app | ? | Development Apps | Status unknown |
| ax-dev-db | ? | Development DB | Status unknown |
| ax-prod-app | ? | Production Apps | Status unknown |
| ax-prod-db | ? | Production DB | Status unknown |

### Observability Stack on ax-sas-tools (121.40.188.25)

**Location:** `/opt/axinova/observability`

**Running Services:**
- **Portainer CE 2.20.3** - https://portainer.axinova-internal.xyz
  - Container management UI
  - Currently managing local endpoint only (endpoint_id: 1)
- **Grafana 11.2.0** - https://grafana.axinova-internal.xyz
  - Metrics visualization and dashboards
  - Current token issues (401 Unauthorized)
- **Prometheus v2.55.0** - https://prometheus.axinova-internal.xyz
  - Metrics collection and storage
  - Currently scraping: localhost:9090 (self) and node exporters on 172.18.80.x
- **Loki 2.9.8** - Running but not exposed via Traefik
  - Log aggregation service
  - Needs to be configured as Grafana datasource

---

## Issues Identified

### 1. ⚠️ Duplicate Observability Stack

**Problem:**
- ax-sas-tools has a full observability stack (Grafana, Prometheus, Loki, Portainer)
- ax-tools was expected to be the centralized observability hub
- ax-tools SSH connection failed (needs investigation)

**Possible Causes:**
1. **Intentional Design:** Each cluster/environment has its own observability stack
2. **Migration in Progress:** Moving from ax-tools to ax-sas-tools
3. **Redundancy:** Both stacks serve different purposes
4. **Misconfiguration:** Duplicate deployment

**Recommendation:**
- **Investigate ax-tools status** - Why is SSH failing?
- **Document intended architecture** - Should there be one central observability stack or distributed?
- **If centralized:** Migrate everything to one location
- **If distributed:** Document the purpose of each stack and configure federation

### 2. 🔴 Grafana Authentication Issue

**Problem:**
- API token returns 401 Unauthorized
- Token was created via internal container API (http://172.26.17.2:3000)
- Now accessing via Traefik (https://grafana.axinova-internal.xyz)

**Root Cause:**
- Token might be for a different Grafana instance
- Token might have expired
- Token permissions might be insufficient

**Recommendation:**
- Create new API token via Traefik-exposed endpoint
- Use service account tokens (Grafana 9+) instead of API keys
- Document token creation process

### 3. ⚠️ Loki Not Configured in Grafana

**Problem:**
- Loki 2.9.8 is running but not configured as a datasource in Grafana
- Log aggregation is not accessible via UI

**Recommendation:**
- Add Loki as a datasource in Grafana
- Configure log shipping from all services to Loki
- Create log viewing dashboards

### 4. ❌ Portainer Agents Not Deployed

**Problem:**
- Portainer only managing local endpoint (ax-sas-tools itself)
- Other machines (ax-dev-app, ax-dev-db, ax-prod-app, ax-prod-db) not connected

**Recommendation:**
- Deploy Portainer agent on all machines
- Configure Portainer to manage all endpoints
- Enable centralized container management

### 5. ⚠️ Limited Monitoring Dashboards

**Problem:**
- No comprehensive monitoring dashboards found
- No machine-specific dashboards
- No service-level dashboards

**Recommendation:**
- Create host monitoring dashboards for all 5 machines
- Create service-specific dashboards for each tech stack
- Follow observability best practices

---

## Recommended Architecture

### Option A: Centralized Observability (Recommended)

```
                   ┌─────────────────────────┐
                   │      ax-tools           │
                   │  (Observability Hub)    │
                   ├─────────────────────────┤
                   │ • Grafana (UI)          │
                   │ • Prometheus (Metrics)  │
                   │ • Loki (Logs)           │
                   │ • Portainer (Containers)│
                   └────────┬────────────────┘
                            │
         ┌──────────────────┼──────────────────┬──────────────────┐
         │                  │                  │                  │
    ┌────▼─────┐      ┌────▼─────┐      ┌────▼─────┐      ┌────▼─────┐
    │ax-sas-   │      │ax-dev-   │      │ax-dev-   │      │ ax-prod- │
    │tools     │      │app       │      │db        │      │ app/db   │
    ├──────────┤      ├──────────┤      ├──────────┤      ├──────────┤
    │•Agent    │      │•Agent    │      │•Agent    │      │•Agent    │
    │•Exporter │      │•Exporter │      │•Exporter │      │•Exporter │
    │•Promtail │      │•Promtail │      │•Promtail │      │•Promtail │
    └──────────┘      └──────────┘      └──────────┘      └──────────┘
```

**Benefits:**
- Single source of truth for all monitoring
- Centralized dashboard management
- Easier to maintain and upgrade
- Lower resource usage

**Implementation Steps:**
1. Fix ax-tools SSH access
2. Deploy agents to all machines
3. Configure metric/log collection
4. Migrate dashboards from ax-sas-tools to ax-tools
5. Deprecate ax-sas-tools observability stack (or keep as backup)

### Option B: Distributed Observability

```
┌─────────────────────┐         ┌─────────────────────┐
│  Production Stack   │         │  Development Stack  │
│  (ax-prod-*)        │         │  (ax-dev-*)         │
├─────────────────────┤         ├─────────────────────┤
│ • Grafana           │         │ • Grafana           │
│ • Prometheus        │         │ • Prometheus        │
│ • Loki              │         │ • Loki              │
└─────────────────────┘         └─────────────────────┘
          │                               │
    ┌─────┴─────┐                   ┌────┴─────┐
    │ ax-prod-  │                   │ ax-dev-  │
    │ app/db    │                   │ app/db   │
    └───────────┘                   └──────────┘

    ┌─────────────────────┐
    │  SaaS Tools Stack   │
    │  (ax-sas-tools)     │
    ├─────────────────────┤
    │ • Grafana           │
    │ • Prometheus        │
    │ • Loki              │
    │ • Portainer (Global)│
    └─────────────────────┘
```

**Benefits:**
- Environment isolation
- Better for multi-tenant setups
- Failure isolation

**Drawbacks:**
- More complex to maintain
- Higher resource usage
- Harder to get unified view

---

## Immediate Action Items

### Priority 1: Fix Critical Issues

1. **Investigate ax-tools Connectivity**
   ```bash
   # Test SSH access
   ssh -v root@121.40.188.18

   # Check if machine is reachable
   ping -c 3 121.40.188.18

   # Try alternative ports
   nmap -p 22,2222 121.40.188.18
   ```

2. **Create New Grafana Token**
   ```bash
   # Via Traefik-exposed endpoint
   curl -k -X POST https://grafana.axinova-internal.xyz/api/auth/keys \
     -u "admin:YOUR_PASSWORD" \
     -H "Content-Type: application/json" \
     -d '{"name":"mcp-server-v2","role":"Admin"}'

   # Update MCP server .env
   APP_GRAFANA__TOKEN=new_token_here
   ```

3. **Configure Loki Datasource**
   ```bash
   # Via MCP or curl
   curl -k -X POST https://grafana.axinova-internal.xyz/api/datasources \
     -H "Authorization: Bearer $GRAFANA_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Loki",
       "type": "loki",
       "url": "http://loki:3100",
       "access": "proxy",
       "isDefault": false
     }'
   ```

### Priority 2: Deploy Portainer Agents

**For each machine (ax-dev-app, ax-dev-db, ax-prod-app, ax-prod-db, ax-sas-tools):**

1. **Install Portainer Agent**
   ```bash
   docker run -d \
     -p 9001:9001 \
     --name portainer_agent \
     --restart=always \
     -v /var/run/docker.sock:/var/run/docker.sock \
     -v /var/lib/docker/volumes:/var/lib/docker/volumes \
     portainer/agent:latest
   ```

2. **Add Endpoint in Portainer**
   - Go to Portainer UI → Environments → Add Environment
   - Choose "Agent"
   - Enter machine IP and port 9001
   - Name it appropriately (e.g., "ax-dev-app")

### Priority 3: Set Up Logging

1. **Install Promtail on All Machines**
   ```yaml
   # docker-compose.yml snippet
   promtail:
     image: grafana/promtail:2.9.8
     volumes:
       - /var/log:/var/log
       - ./promtail-config.yml:/etc/promtail/config.yml
     command: -config.file=/etc/promtail/config.yml
   ```

2. **Configure Promtail**
   ```yaml
   # promtail-config.yml
   server:
     http_listen_port: 9080

   positions:
     filename: /tmp/positions.yaml

   clients:
     - url: http://loki:3100/loki/api/v1/push

   scrape_configs:
     - job_name: docker
       static_configs:
         - targets:
             - localhost
           labels:
             job: docker
             host: ${HOSTNAME}
             __path__: /var/lib/docker/containers/*/*.log
   ```

3. **Create Log Viewing Dashboard**
   - Use Grafana's Loki datasource
   - Add filters for: service, container, severity, host
   - Include log volume metrics
   - Add log pattern detection

### Priority 4: Create Monitoring Dashboards

**For Each Machine:**

1. **Host Metrics Dashboard**
   - CPU usage (per core and total)
   - Memory usage (used, cached, available)
   - Disk usage and I/O
   - Network traffic
   - System load

2. **Container Metrics Dashboard**
   - Container CPU usage
   - Container memory usage
   - Container network I/O
   - Container restart count
   - Container health status

3. **Service-Specific Dashboards**
   - Go applications: goroutines, GC stats, memory allocations
   - PostgreSQL: connections, queries/sec, cache hit rate
   - Redis: memory usage, keyspace, commands/sec
   - Web services: request rate, latency, error rate

---

## Logging Dashboard Specification

### Log Viewer Dashboard Requirements

**Filters:**
1. **Time Range** - Last 5m, 15m, 1h, 3h, 6h, 12h, 24h, 7d, custom
2. **Host/Machine** - ax-tools, ax-sas-tools, ax-dev-app, ax-dev-db, ax-prod-app, ax-prod-db
3. **Service** - All services running on the machines
4. **Container** - Individual Docker containers
5. **Log Level** - DEBUG, INFO, WARNING, ERROR, FATAL
6. **Search** - Full-text search in log messages

**Panels:**
1. **Log Volume Over Time** - Bar chart showing log count per time bucket
2. **Log Level Distribution** - Pie chart of log levels
3. **Top Log Sources** - Table of top services/containers by log count
4. **Live Log Stream** - Real-time log tail with syntax highlighting
5. **Log Pattern Detection** - Common log patterns and anomalies
6. **Error Rate** - Time series of error log count

**Example LogQL Queries:**
```logql
# All logs from a specific host
{host="ax-sas-tools"}

# Error logs from all services
{job="docker"} |= "level=error"

# Logs from specific container
{container_name="axinova-mcp-server"}

# Go application logs with JSON parsing
{job="docker"} | json | level="error"

# Pattern matching
{job="docker"} |~ "error|exception|fatal"
```

---

## Monitoring Dashboard Specification

### Host Monitoring Dashboard (per machine)

**Panels:**

1. **System Overview**
   - Uptime
   - OS Version
   - Kernel Version
   - Architecture

2. **CPU Metrics**
   ```promql
   # CPU Usage %
   100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m])) * 100)

   # CPU per core
   irate(node_cpu_seconds_total[5m]) * 100

   # Load Average
   node_load1, node_load5, node_load15
   ```

3. **Memory Metrics**
   ```promql
   # Memory Usage %
   (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100

   # Memory breakdown
   node_memory_MemTotal_bytes
   node_memory_MemAvailable_bytes
   node_memory_Cached_bytes
   node_memory_Buffers_bytes
   ```

4. **Disk Metrics**
   ```promql
   # Disk Usage %
   (1 - (node_filesystem_avail_bytes / node_filesystem_size_bytes)) * 100

   # Disk I/O
   irate(node_disk_read_bytes_total[5m])
   irate(node_disk_written_bytes_total[5m])

   # Disk IOPS
   irate(node_disk_reads_completed_total[5m])
   irate(node_disk_writes_completed_total[5m])
   ```

5. **Network Metrics**
   ```promql
   # Network Traffic
   irate(node_network_receive_bytes_total[5m])
   irate(node_network_transmit_bytes_total[5m])

   # Network Errors
   irate(node_network_receive_errs_total[5m])
   irate(node_network_transmit_errs_total[5m])
   ```

### Container Monitoring Dashboard

**Panels:**

1. **Container List** - Table with: Name, Image, Status, Uptime, CPU%, Memory%

2. **Container CPU**
   ```promql
   rate(container_cpu_usage_seconds_total{name!=""}[5m]) * 100
   ```

3. **Container Memory**
   ```promql
   container_memory_usage_bytes{name!=""}
   container_memory_working_set_bytes{name!=""}
   ```

4. **Container Network**
   ```promql
   rate(container_network_receive_bytes_total[5m])
   rate(container_network_transmit_bytes_total[5m])
   ```

5. **Container Restarts**
   ```promql
   increase(container_start_time_seconds[1h])
   ```

---

## Tech Stack Best Practices

### Go Applications

**Metrics to Monitor:**
- Goroutines count
- Memory allocations
- GC pause times
- HTTP request latency (p50, p95, p99)
- HTTP error rate
- Custom business metrics

**Recommended:**
- Use Prometheus client library
- Expose `/metrics` endpoint
- Use structured logging (logrus, zap)
- Add trace IDs to logs

### PostgreSQL

**Metrics to Monitor:**
- Active connections
- Transaction rate
- Query duration
- Cache hit ratio
- Index usage
- Replication lag (if applicable)

**Recommended:**
- Use postgres_exporter
- Monitor slow query log
- Track table bloat
- Monitor vacuum activity

### Redis

**Metrics to Monitor:**
- Memory usage
- Connected clients
- Commands/sec
- Hit rate
- Evicted keys
- Replication lag

**Recommended:**
- Use redis_exporter
- Monitor key expiration
- Track slow commands
- Monitor persistence

---

## Implementation Timeline

### Week 1: Foundation
- [ ] Investigate ax-tools connectivity
- [ ] Document intended architecture
- [ ] Fix Grafana token issues
- [ ] Configure Loki datasource
- [ ] Deploy Portainer agents to 2 machines

### Week 2: Logging
- [ ] Deploy Promtail to all machines
- [ ] Configure log shipping
- [ ] Create log viewing dashboard
- [ ] Test log filtering and search
- [ ] Deploy remaining Portainer agents

### Week 3: Monitoring
- [ ] Create host monitoring dashboards (5 machines)
- [ ] Create container monitoring dashboard
- [ ] Create service-specific dashboards
- [ ] Set up alerts for critical metrics
- [ ] Document dashboard usage

### Week 4: Polish & Training
- [ ] Optimize queries and performance
- [ ] Add custom dashboards per team needs
- [ ] Create runbooks for common issues
- [ ] Train team on dashboard usage
- [ ] Document maintenance procedures

---

## Using MCP Server for Implementation

The newly deployed MCP server can be used for these tasks!

**Example: Configure Loki Datasource**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "grafana_create_datasource",
    "arguments": {
      "name": "Loki",
      "type": "loki",
      "url": "http://loki:3100"
    }
  }
}
```

**Example: Create Dashboard**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "grafana_create_dashboard",
    "arguments": {
      "title": "Host Monitoring - ax-sas-tools",
      "dashboard": {
        // Dashboard JSON definition
      }
    }
  }
}
```

**Example: Check Container Status**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "portainer_list_containers",
    "arguments": {
      "endpoint_id": 1
    }
  }
}
```

---

## Maintenance & Best Practices

### Dashboards
- Use folders to organize dashboards by category
- Add descriptions to panels
- Use variables for dynamic filtering
- Set appropriate time ranges
- Add links between related dashboards

### Alerts
- Set meaningful thresholds based on SLOs
- Use alert grouping to reduce noise
- Configure notification channels (email, Slack, PagerDuty)
- Add runbook links to alerts
- Test alerts regularly

### Data Retention
- Prometheus: 15 days (adjustable)
- Loki: 30 days (adjustable)
- Consider long-term storage solutions for compliance

### Security
- Rotate API tokens regularly
- Use RBAC in Grafana
- Limit Portainer agent access
- Secure Prometheus/Loki endpoints
- Enable audit logging

---

## Conclusion

The current setup has a solid foundation with Grafana, Prometheus, and Loki running on ax-sas-tools. The next steps are to:

1. Clarify the intended architecture (centralized vs distributed)
2. Fix authentication issues
3. Configure Loki for log aggregation
4. Deploy Portainer agents to all machines
5. Create comprehensive monitoring and logging dashboards

With the MCP server now operational, many of these tasks can be automated and executed programmatically, making infrastructure management more efficient and reproducible.

---

**For Questions or Assistance:**
- MCP Server Documentation: `/opt/axinova-mcp-server/`
- Grafana UI: https://grafana.axinova-internal.xyz
- Portainer UI: https://portainer.axinova-internal.xyz
- Prometheus UI: https://prometheus.axinova-internal.xyz
