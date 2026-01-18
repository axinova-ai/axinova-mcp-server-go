package prometheus

import (
	"context"
	"fmt"
	"time"

	"github.com/axinova-ai/axinova-mcp-server-go/internal/mcp"
)

// RegisterTools registers all Prometheus tools with the MCP server
func RegisterTools(server *mcp.Server, client *Client) {
	// Query instant
	server.RegisterTool(mcp.Tool{
		Name:        "prometheus_query",
		Description: "Execute an instant Prometheus query",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"query": {
					Type:        "string",
					Description: "PromQL query expression",
				},
			},
			Required: []string{"query"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		query, ok := args["query"].(string)
		if !ok {
			return nil, fmt.Errorf("query is required")
		}

		result, err := client.Query(ctx, query, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %w", err)
		}
		return result, nil
	})

	// Query range
	server.RegisterTool(mcp.Tool{
		Name:        "prometheus_query_range",
		Description: "Execute a Prometheus range query over a time period",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"query": {
					Type:        "string",
					Description: "PromQL query expression",
				},
				"start": {
					Type:        "string",
					Description: "Start time (RFC3339 format or relative like '1h' ago)",
				},
				"end": {
					Type:        "string",
					Description: "End time (RFC3339 format or 'now', default: now)",
				},
				"step": {
					Type:        "string",
					Description: "Query resolution step (e.g., '15s', '1m', default: 1m)",
				},
			},
			Required: []string{"query", "start"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		query, ok := args["query"].(string)
		if !ok {
			return nil, fmt.Errorf("query is required")
		}

		startStr, ok := args["start"].(string)
		if !ok {
			return nil, fmt.Errorf("start is required")
		}

		// Parse start time
		start, err := parseTimeOrRelative(startStr)
		if err != nil {
			return nil, fmt.Errorf("invalid start time: %w", err)
		}

		// Parse end time (default: now)
		end := time.Now()
		if endStr, ok := args["end"].(string); ok && endStr != "now" {
			end, err = parseTimeOrRelative(endStr)
			if err != nil {
				return nil, fmt.Errorf("invalid end time: %w", err)
			}
		}

		// Step (default: 1m)
		step := "1m"
		if s, ok := args["step"].(string); ok {
			step = s
		}

		result, err := client.QueryRange(ctx, query, start, end, step)
		if err != nil {
			return nil, fmt.Errorf("failed to execute range query: %w", err)
		}
		return result, nil
	})

	// List label names
	server.RegisterTool(mcp.Tool{
		Name:        "prometheus_list_label_names",
		Description: "Get all Prometheus label names",
		InputSchema: mcp.InputSchema{
			Type: "object",
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		labels, err := client.LabelNames(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get label names: %w", err)
		}
		return labels, nil
	})

	// List label values
	server.RegisterTool(mcp.Tool{
		Name:        "prometheus_list_label_values",
		Description: "Get all values for a specific Prometheus label",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"label": {
					Type:        "string",
					Description: "Label name (e.g., 'job', 'instance', '__name__')",
				},
			},
			Required: []string{"label"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		label, ok := args["label"].(string)
		if !ok {
			return nil, fmt.Errorf("label is required")
		}

		values, err := client.LabelValues(ctx, label)
		if err != nil {
			return nil, fmt.Errorf("failed to get label values: %w", err)
		}
		return values, nil
	})

	// Find series
	server.RegisterTool(mcp.Tool{
		Name:        "prometheus_find_series",
		Description: "Find time series by label matchers",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"match": {
					Type:        "string",
					Description: "Series selector (e.g., 'up{job=\"prometheus\"}')",
				},
				"lookback": {
					Type:        "string",
					Description: "How far back to look (e.g., '1h', '24h', default: 1h)",
				},
			},
			Required: []string{"match"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		match, ok := args["match"].(string)
		if !ok {
			return nil, fmt.Errorf("match is required")
		}

		// Lookback period
		lookback := "1h"
		if lb, ok := args["lookback"].(string); ok {
			lookback = lb
		}

		duration, err := time.ParseDuration(lookback)
		if err != nil {
			return nil, fmt.Errorf("invalid lookback duration: %w", err)
		}

		end := time.Now()
		start := end.Add(-duration)

		series, err := client.Series(ctx, []string{match}, start, end)
		if err != nil {
			return nil, fmt.Errorf("failed to find series: %w", err)
		}
		return series, nil
	})

	// List targets
	server.RegisterTool(mcp.Tool{
		Name:        "prometheus_list_targets",
		Description: "Get all Prometheus scrape targets and their health status",
		InputSchema: mcp.InputSchema{
			Type: "object",
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		targets, err := client.Targets(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get targets: %w", err)
		}
		return targets, nil
	})

	// Get metadata
	server.RegisterTool(mcp.Tool{
		Name:        "prometheus_get_metadata",
		Description: "Get metric metadata (HELP and TYPE information)",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"metric": {
					Type:        "string",
					Description: "Metric name (optional, returns all if not specified)",
				},
			},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		metric := ""
		if m, ok := args["metric"].(string); ok {
			metric = m
		}

		metadata, err := client.Metadata(ctx, metric)
		if err != nil {
			return nil, fmt.Errorf("failed to get metadata: %w", err)
		}
		return metadata, nil
	})
}

// parseTimeOrRelative parses a time string or relative duration
func parseTimeOrRelative(s string) (time.Time, error) {
	// Try parsing as RFC3339
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}

	// Try parsing as relative duration (e.g., "1h", "30m")
	if duration, err := time.ParseDuration(s); err == nil {
		return time.Now().Add(-duration), nil
	}

	return time.Time{}, fmt.Errorf("invalid time format: %s (use RFC3339 or duration like '1h')", s)
}
