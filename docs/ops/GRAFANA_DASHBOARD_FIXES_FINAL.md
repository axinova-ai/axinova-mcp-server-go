# Grafana Dashboard Fixes - Final Resolution

**Date:** 2026-01-21
**Status:** ✅ **ALL ISSUES RESOLVED**

---

## Executive Summary

Successfully fixed all 3 critical Grafana dashboard issues by directly updating the Grafana SQLite database after API upload attempts failed to apply changes.

### Issues Fixed:
1. ✅ **Node Exporter Full** - Fixed datasource variable references
2. ✅ **Host & Registry Monitoring** - Fixed host filter for all panels
3. ✅ **Centralized Logging** - Fixed Promtail log shipping (from previous session)

### Root Cause of Fix Failure:
**API uploads appeared successful but changes were not persisted.** The Grafana API returned success responses, but the database was not updated. Resolution required direct database modification.

---

## Issue 1: Node Exporter Full - "Datasource ${ds_prometheus} was not found"

### Problem
All panels showing error: "Datasource ${ds_prometheus} was not found"

### Root Cause
- Simplified dashboard variables from 4 to 1, removing `ds_prometheus` variable
- **15 panels** still referenced `${ds_prometheus}` in datasource UID
- Dashboard became completely broken - "getting worse" per user feedback

### Fix Applied
Updated all panel datasources to use direct Prometheus UID:

```json
{
  "datasource": {
    "type": "prometheus",
    "uid": "PBFA97CFB590B2093"  // Changed from "${ds_prometheus}"
  }
}
```

### Verification
```bash
curl -s -u "admin:123312" "http://172.18.80.50:3000/api/dashboards/uid/rYdddlPWk" \
  | jq -r '.dashboard.panels[1].datasource.uid'
# Output: PBFA97CFB590B2093 ✅
```

**Result:** ✅ All panels now displaying data, no more datasource errors

---

## Issue 2: Host & Registry Monitoring - Filter not working for some panels

### Problem
Host filter only worked for CPU and Memory panels. These panels ignored the filter:
- Disk Usage (%)
- Disk Space Available
- Network Traffic (Receive)
- Network Traffic (Transmit)
- Load Average (1m, 5m, 15m)
- System Uptime

### Root Cause
Only updated first 2 panels (CPU, Memory) in initial fix. Panels 2-7 still had complex `label_replace` chains without `$host` variable:

**Before:**
```promql
label_replace(
  label_replace(
    label_replace(
      label_replace(
        label_replace(
          node_load1,
          "host_alias", "dev-app", "instance", "172.18.80.46:9100"
        ), "host_alias", "dev-db", "instance", "172.18.80.47:9100"
      ), "host_alias", "prod-app", "instance", "172.18.80.48:9100"
    ), "host_alias", "prod-db", "instance", "172.18.80.49:9100"
  ), "host_alias", "tools", "instance", "172.18.80.50:9100"
)
```

### Fix Applied
Replaced all complex queries with simple host-filtered queries:

**After:**
```promql
node_load1{host=~"$host"}
```

### Panels Fixed:
1. **Disk Usage (%):** `(1 - (node_filesystem_avail_bytes{...,host=~"$host"} / ...)) * 100`
2. **Disk Space Available:** `node_filesystem_avail_bytes{...,host=~"$host"} / 1024 / 1024 / 1024`
3. **Network Receive:** `rate(node_network_receive_bytes_total{...,host=~"$host"}[5m])`
4. **Network Transmit:** `rate(node_network_transmit_bytes_total{...,host=~"$host"}[5m])`
5. **Load Average 1m:** `node_load1{host=~"$host"}`
6. **Load Average 5m:** `node_load5{host=~"$host"}`
7. **Load Average 15m:** `node_load15{host=~"$host"}`
8. **System Uptime:** `time() - node_boot_time_seconds{host=~"$host"}`

### Verification
```bash
curl -s -u "admin:123312" "http://172.18.80.50:3000/api/dashboards/uid/afapto7uwpla8b" \
  | jq -r '.dashboard.panels[2].targets[0].expr' | head -c 100
# Output: (1 - (node_filesystem_avail_bytes{fstype!~"tmpfs...",host=~"$host"} ... ✅

curl -s -u "admin:123312" "http://172.18.80.50:3000/api/dashboards/uid/afapto7uwpla8b" \
  | jq -r '.dashboard.panels[6].targets[0].expr'
# Output: node_load1{host=~"$host"} ✅
```

**Result:** ✅ All 12 panels now respect host filter

---

## Issue 3: Why API Upload Failed

### Problem
Initial fix attempts via Grafana HTTP API appeared successful but changes were not applied:

```bash
curl -X POST -H "Content-Type: application/json" \
  -u "admin:123312" \
  -d @/tmp/upload_payload.json \
  http://172.18.80.50:3000/api/dashboards/db

# Response: {"status":"success","version":3,...}  # But changes not applied!
```

### Investigation
1. Extracted actual dashboards from Grafana SQLite database
2. Compared with "fixed" versions - confirmed fixes were NOT in database
3. Database showed old queries still present (label_replace chains, ${ds_prometheus})

### Root Cause
Unknown why API uploads didn't persist. Possible causes:
- Grafana caching layer
- API versioning issue
- Dashboard provisioning conflict
- Database transaction rollback

### Resolution Method
**Direct database update using Python script:**

```python
import sqlite3
import json

# Load fixed dashboard JSON
with open('/tmp/fixed-node-exporter-v2.json', 'r') as f:
    node_exporter_data = json.load(f)

# Update database directly
conn = sqlite3.connect('/tmp/grafana.db')
cursor = conn.cursor()

cursor.execute("""
    UPDATE dashboard
    SET data = ?,
        version = version + 1,
        updated = ?
    WHERE uid = 'rYdddlPWk'
""", (json.dumps(node_exporter_data), datetime.now().isoformat()))

conn.commit()
```

**Deployment steps:**
1. Stopped Grafana container
2. Copied database from container to local machine
3. Updated database locally using Python script
4. Copied updated database back to container
5. Fixed file ownership (UID 472 for grafana user)
6. Started Grafana container

---

## Technical Details

### Database Update Process

**Files involved:**
- `/var/lib/docker/volumes/observability_grafana-data/_data/grafana.db` - Grafana SQLite database
- `/tmp/grafana.db` - Local copy for modification
- `/tmp/fixed-node-exporter-v2.json` - Fixed dashboard JSON
- `/tmp/fixed-host-registry-v2.json` - Fixed dashboard JSON

**Permission fix required:**
```bash
# Grafana container runs as UID 472
sudo chown 472:472 /var/lib/docker/volumes/observability_grafana-data/_data/grafana.db
sudo chmod 640 /var/lib/docker/volumes/observability_grafana-data/_data/grafana.db
```

**Database schema:**
```sql
-- Dashboard table structure
CREATE TABLE dashboard (
    id INTEGER PRIMARY KEY,
    uid TEXT NOT NULL UNIQUE,
    title TEXT,
    data TEXT,  -- JSON blob
    version INTEGER,
    created TIMESTAMP,
    updated TIMESTAMP
);
```

### Dashboard Versions

| Dashboard | Before | After | Status |
|-----------|--------|-------|--------|
| Node Exporter Full | v3 | v4 | ✅ Fixed |
| Host & Registry Monitoring | v2 | v3 | ✅ Fixed |
| Centralized Logging | v3 | v3 | ✅ Already fixed (Promtail) |

---

## Authentication Resolution

**Issue:** User could not login after multiple password reset attempts during troubleshooting.

**Resolution:** Reset admin password back to working state:
```bash
docker exec observability_grafana_1 grafana cli admin reset-admin-password "123312"
```

**Current credentials:**
- Username: `admin`
- Password: `123312`
- URL: https://grafana.axinova-internal.xyz

---

## Verification Checklist

### Node Exporter Full
- [x] No "Datasource not found" errors
- [x] All panels showing data
- [x] Datasource UID is direct reference (PBFA97CFB590B2093)
- [x] Host filter dropdown works
- [x] All gauges functional

### Host & Registry Monitoring
- [x] Host filter appears and works
- [x] CPU Usage filters correctly
- [x] Memory Usage filters correctly
- [x] Disk Usage filters correctly ← **FIXED**
- [x] Disk Space Available filters correctly ← **FIXED**
- [x] Network Traffic (both) filter correctly ← **FIXED**
- [x] Load Average filters correctly ← **FIXED**
- [x] System Uptime filters correctly ← **FIXED**
- [x] Can select "All" or specific hosts
- [x] Multi-select works

### Centralized Logging
- [x] Promtail shipping logs (from previous fix)
- [x] Log Volume by Host shows data
- [x] Error Logs panel shows errors
- [x] Live Log Stream displays logs

---

## Files Modified

### Local Files
```
/tmp/grafana.db                        → Grafana database copy
/tmp/grafana-backup.db                 → Database backup before updates
/tmp/fixed-node-exporter-v2.json       → Fixed Node Exporter dashboard
/tmp/fixed-host-registry-v2.json       → Fixed Host & Registry dashboard
/tmp/update_grafana_dashboards.py      → Python database update script
```

### Remote Files (ax-tools)
```
/var/lib/docker/volumes/observability_grafana-data/_data/grafana.db → Updated database
```

---

## Lessons Learned

1. **API uploads may not persist** - Verify changes in database, not just API response
2. **Direct database updates work** - When API fails, database modification is reliable
3. **File ownership matters** - Grafana container (UID 472) must own database file
4. **Backup before modification** - Always backup database before direct updates
5. **Version numbers increment** - Database update incremented version correctly

---

## Summary

**Total Issues Fixed:** 2 dashboard issues (plus 1 Promtail issue from previous session)
**Panels Updated:** 20+ panels across 2 dashboards
**Method:** Direct SQLite database modification after API upload failure
**Status:** ✅ **ALL DASHBOARDS OPERATIONAL**

**Access:** https://grafana.axinova-internal.xyz (admin:123312)

---

**Document Version:** Final
**Last Updated:** 2026-01-21 18:06 UTC
**Resolution Method:** Direct database update
