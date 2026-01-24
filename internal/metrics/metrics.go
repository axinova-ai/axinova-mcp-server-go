package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	startTime = time.Now()

	// Server uptime
	UptimeSeconds = promauto.NewGaugeFunc(
		prometheus.GaugeOpts{
			Name: "mcp_server_uptime_seconds",
			Help: "Server uptime in seconds",
		},
		func() float64 {
			return time.Since(startTime).Seconds()
		},
	)

	// Number of registered tools
	RegisteredTools = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mcp_tools_registered_total",
		Help: "Total number of registered MCP tools",
	})

	// Number of registered resources
	RegisteredResources = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mcp_resources_registered_total",
		Help: "Total number of registered MCP resources",
	})

	// Total RPC requests handled
	RPCRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mcp_rpc_requests_total",
			Help: "Total number of RPC requests handled",
		},
		[]string{"method", "transport"}, // transport: stdio or http
	)

	// RPC request duration
	RPCRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "mcp_rpc_request_duration_seconds",
			Help:    "RPC request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "transport"},
	)

	// RPC errors
	RPCErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "mcp_rpc_errors_total",
			Help: "Total number of RPC errors",
		},
		[]string{"method", "error_code", "transport"},
	)

	// Active connections (for HTTP mode)
	ActiveConnections = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mcp_http_active_connections",
		Help: "Number of active HTTP connections",
	})
)

// RecordToolsRegistered updates the count of registered tools
func RecordToolsRegistered(count int) {
	RegisteredTools.Set(float64(count))
}

// RecordResourcesRegistered updates the count of registered resources
func RecordResourcesRegistered(count int) {
	RegisteredResources.Set(float64(count))
}

// RecordRPCRequest records a completed RPC request
func RecordRPCRequest(method, transport string, duration time.Duration, errCode string) {
	RPCRequestsTotal.WithLabelValues(method, transport).Inc()
	RPCRequestDuration.WithLabelValues(method, transport).Observe(duration.Seconds())

	if errCode != "" {
		RPCErrorsTotal.WithLabelValues(method, errCode, transport).Inc()
	}
}
