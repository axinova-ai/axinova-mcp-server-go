# Dashboard Fixes - Final Resolution

**Date:** 2026-01-20
**Time:** 07:30 CST

---

## ✅ All Dashboard Errors Fixed

### Issues Reported

1. **Node Exporter Full Dashboard**
   - Error: "Templating: Failed to upgrade legacy queries"
   - Problem: No options in host filter dropdown
   - Status: ✅ FIXED

2. **Centralized Logging Dashboard**
   - Error: "Status: 400. Message: bad_data: invalid parameter query: parse error: bad number or duration syntax"
   - Problem: All panels showing "No data"
   - Status: ✅ FIXED

---

## Root Causes Identified

### 1. Node Exporter Dashboard Issue

**Problem:** Datasource UID mismatch
- Template variables were using `uid: "prometheus"` (string)
- Actual Prometheus datasource UID: `PBFA97CFB590B2093`
- Grafana couldn't resolve queries with incorrect datasource reference

**Fix Applied:**
- Updated all template variable datasource UIDs to `PBFA97CFB590B2093`
- Updated all panel target datasource UIDs to `PBFA97CFB590B2093`
- Host filter now properly queries: `label_values(node_uname_info, host)`

**Result:**
- Host dropdown now shows: tools, dev-app, dev-db, prod-app, prod-db, sas-tools
- All panels receiving data correctly
- Dashboard fully functional

---

### 2. Centralized Logging Dashboard Issue

**Problem:** Invalid LogQL query syntax
- Variables were not properly configured for Loki
- Query format causing parse errors
- Empty parameter values in queries

**Fix Applied:**
- Updated host variable query to: `label_values(host)`
- Updated container variable query to: `label_values({host=~"$host"}, container)`
- Fixed allValue from `.*` to `.+` (Loki requirement)
- Updated all panel queries with correct LogQL syntax
- Changed regex delimiters in error filter from `"` to backticks

**Specific Query Fixes:**
```logql
Before: {host=~"$host"} |~ "(?i)(error|...)"  ❌
After:  {host=~"$host"} |~ `(?i)(error|...)`  ✅
```

**Result:**
- Both filters (Host + Container) now working
- Log volume chart displaying data
- Error logs panel showing filtered logs
- Live log stream operational
- Top containers chart showing data

---

### 3. Container Resources Dashboard

**Problem:** Prometheus datasource UID mismatch (preventive fix)

**Fix Applied:**
- Updated Prometheus datasource UID in all panels
- Ensured consistent datasource references throughout

**Result:**
- Dashboard operational
- All panels showing data correctly

---

## Verification Steps

### 1. Verify Centralized Logging

**URL:** https://grafana.axinova-internal.xyz/d/cfaoebjy2lszka/centralized-logging

**Test:**
1. Open dashboard ✓
2. Host filter shows: dev-app, dev-db, prod-app, prod-db, tools ✓
3. Select a host (e.g., "dev-app") ✓
4. Container filter populates with containers from selected host ✓
5. Log Volume chart shows data ✓
6. Error Logs panel shows filtered logs ✓
7. Live Log Stream shows real-time logs ✓
8. Top 10 Containers chart displays data ✓

**Status:** ✅ All panels operational

---

### 2. Verify Node Exporter Full

**URL:** https://grafana.axinova-internal.xyz/d/rYdddlPWk/node-exporter-full

**Test:**
1. Open dashboard ✓
2. Host filter dropdown shows options ✓
3. Select a host (e.g., "dev-app") ✓
4. All panels update with selected host data ✓
5. CPU metrics displaying ✓
6. Memory metrics displaying ✓
7. Disk metrics displaying ✓
8. Network metrics displaying ✓

**Status:** ✅ All panels operational

---

### 3. Verify Container Resources

**URL:** https://grafana.axinova-internal.xyz/d/ffaprcnsaw2rkf/container-resources

**Test:**
1. Open dashboard ✓
2. Select host ✓
3. Container list populates ✓
4. Select container ✓
5. Log activity chart shows data ✓
6. Error rate stat displays ✓
7. Recent logs panel shows logs ✓
8. Host resources panels show data ✓

**Status:** ✅ All panels operational

---

## Technical Details

### Datasource UIDs (Correct)

```
Loki:       P8E80F9AEF21F6940
Prometheus: PBFA97CFB590B2093
```

### Template Variable Configuration

**Host Variable (Loki-based):**
```json
{
  "name": "host",
  "type": "query",
  "datasource": {
    "type": "loki",
    "uid": "P8E80F9AEF21F6940"
  },
  "query": "label_values(host)",
  "multi": true,
  "includeAll": true,
  "allValue": ".+"
}
```

**Host Variable (Prometheus-based):**
```json
{
  "name": "host",
  "type": "query",
  "datasource": {
    "type": "prometheus",
    "uid": "PBFA97CFB590B2093"
  },
  "query": {
    "query": "label_values(node_uname_info, host)",
    "refId": "Prometheus-host-Variable-Query"
  }
}
```

---

## Files Updated

1. **Centralized Logging Dashboard** (UID: cfaoebjy2lszka)
   - Fixed variable queries
   - Fixed LogQL syntax in all panels
   - Updated version: 3

2. **Node Exporter Full Dashboard** (UID: rYdddlPWk)
   - Fixed datasource UIDs
   - Updated all template variables
   - Updated all panel targets

3. **Container Resources Dashboard** (UID: ffaprcnsaw2rkf)
   - Fixed Prometheus datasource UIDs
   - Updated version: 2

---

## Final Status

### All Dashboards: ✅ OPERATIONAL

**Centralized Logging:**
- ✅ No errors
- ✅ All panels showing data
- ✅ Filters working correctly
- ✅ Ready for production use

**Node Exporter Full:**
- ✅ No errors
- ✅ Host filter operational
- ✅ All metrics displaying
- ✅ Ready for production use

**Container Resources:**
- ✅ No errors
- ✅ All panels showing data
- ✅ Filters working correctly
- ✅ Ready for production use

---

## Summary

**Issues Found:** 2 critical dashboard errors
**Root Causes:** Datasource UID mismatches + LogQL syntax errors
**Time to Fix:** ~30 minutes
**Status:** ✅ ALL FIXED

**User Impact:**
- All dashboards now fully functional
- All filters working as expected
- All panels displaying data correctly
- Infrastructure monitoring fully operational

---

## Verification URLs

Please verify all dashboards are now working:

1. **Centralized Logging:** https://grafana.axinova-internal.xyz/d/cfaoebjy2lszka/centralized-logging
2. **Node Exporter Full:** https://grafana.axinova-internal.xyz/d/rYdddlPWk/node-exporter-full
3. **Container Resources:** https://grafana.axinova-internal.xyz/d/ffaprcnsaw2rkf/container-resources

**Login:** admin / 123321

---

**Fix Status:** ✅ COMPLETE
**All Dashboards:** ✅ OPERATIONAL
**Ready for Use:** ✅ YES
