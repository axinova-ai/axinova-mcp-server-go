package grafana

import (
	"context"
	"fmt"

	"github.com/axinova-ai/axinova-mcp-server-go/internal/mcp"
)

// RegisterTools registers all Grafana tools with the MCP server
func RegisterTools(server *mcp.Server, client *Client) {
	// List dashboards
	server.RegisterTool(mcp.Tool{
		Name:        "grafana_list_dashboards",
		Description: "List all Grafana dashboards",
		InputSchema: mcp.InputSchema{
			Type: "object",
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		dashboards, err := client.ListDashboards(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list dashboards: %w", err)
		}
		return dashboards, nil
	})

	// Get dashboard
	server.RegisterTool(mcp.Tool{
		Name:        "grafana_get_dashboard",
		Description: "Get a specific Grafana dashboard by UID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid": {
					Type:        "string",
					Description: "Dashboard UID",
				},
			},
			Required: []string{"uid"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		uid, ok := args["uid"].(string)
		if !ok {
			return nil, fmt.Errorf("uid is required")
		}

		dashboard, err := client.GetDashboard(ctx, uid)
		if err != nil {
			return nil, fmt.Errorf("failed to get dashboard: %w", err)
		}
		return dashboard, nil
	})

	// Create dashboard
	server.RegisterTool(mcp.Tool{
		Name:        "grafana_create_dashboard",
		Description: "Create a new Grafana dashboard",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"title": {
					Type:        "string",
					Description: "Dashboard title",
				},
				"folder_uid": {
					Type:        "string",
					Description: "Folder UID (optional, default: General)",
				},
				"dashboard": {
					Type:        "object",
					Description: "Dashboard JSON definition (optional, will create basic dashboard if not provided)",
				},
			},
			Required: []string{"title"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		title, ok := args["title"].(string)
		if !ok {
			return nil, fmt.Errorf("title is required")
		}

		folderUID := ""
		if uid, ok := args["folder_uid"].(string); ok {
			folderUID = uid
		}

		// Use provided dashboard or create a basic one
		dashboard := map[string]interface{}{
			"title": title,
			"tags":  []string{"mcp"},
		}
		if customDash, ok := args["dashboard"].(map[string]interface{}); ok {
			dashboard = customDash
			dashboard["title"] = title
		}

		result, err := client.CreateDashboard(ctx, dashboard, folderUID, false)
		if err != nil {
			return nil, fmt.Errorf("failed to create dashboard: %w", err)
		}
		return result, nil
	})

	// Delete dashboard
	server.RegisterTool(mcp.Tool{
		Name:        "grafana_delete_dashboard",
		Description: "Delete a Grafana dashboard by UID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"uid": {
					Type:        "string",
					Description: "Dashboard UID",
				},
			},
			Required: []string{"uid"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		uid, ok := args["uid"].(string)
		if !ok {
			return nil, fmt.Errorf("uid is required")
		}

		if err := client.DeleteDashboard(ctx, uid); err != nil {
			return nil, fmt.Errorf("failed to delete dashboard: %w", err)
		}
		return fmt.Sprintf("Dashboard %s deleted successfully", uid), nil
	})

	// List datasources
	server.RegisterTool(mcp.Tool{
		Name:        "grafana_list_datasources",
		Description: "List all Grafana datasources",
		InputSchema: mcp.InputSchema{
			Type: "object",
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		datasources, err := client.ListDatasources(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list datasources: %w", err)
		}
		return datasources, nil
	})

	// Create datasource
	server.RegisterTool(mcp.Tool{
		Name:        "grafana_create_datasource",
		Description: "Create a new Grafana datasource",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"name": {
					Type:        "string",
					Description: "Datasource name",
				},
				"type": {
					Type:        "string",
					Description: "Datasource type (e.g., prometheus, loki, elasticsearch)",
				},
				"url": {
					Type:        "string",
					Description: "Datasource URL",
				},
				"is_default": {
					Type:        "boolean",
					Description: "Set as default datasource",
				},
			},
			Required: []string{"name", "type", "url"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		name, ok := args["name"].(string)
		if !ok {
			return nil, fmt.Errorf("name is required")
		}

		dsType, ok := args["type"].(string)
		if !ok {
			return nil, fmt.Errorf("type is required")
		}

		url, ok := args["url"].(string)
		if !ok {
			return nil, fmt.Errorf("url is required")
		}

		isDefault := false
		if def, ok := args["is_default"].(bool); ok {
			isDefault = def
		}

		ds := Datasource{
			Name:      name,
			Type:      dsType,
			URL:       url,
			IsDefault: isDefault,
		}

		result, err := client.CreateDatasource(ctx, ds)
		if err != nil {
			return nil, fmt.Errorf("failed to create datasource: %w", err)
		}
		return result, nil
	})

	// Query datasource
	server.RegisterTool(mcp.Tool{
		Name:        "grafana_query_datasource",
		Description: "Query a Grafana datasource (Prometheus, Loki, etc.)",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"datasource_uid": {
					Type:        "string",
					Description: "Datasource UID",
				},
				"query": {
					Type:        "string",
					Description: "Query expression (e.g., PromQL for Prometheus)",
				},
			},
			Required: []string{"datasource_uid", "query"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		dsUID, ok := args["datasource_uid"].(string)
		if !ok {
			return nil, fmt.Errorf("datasource_uid is required")
		}

		query, ok := args["query"].(string)
		if !ok {
			return nil, fmt.Errorf("query is required")
		}

		result, err := client.QueryDatasource(ctx, dsUID, query)
		if err != nil {
			return nil, fmt.Errorf("failed to query datasource: %w", err)
		}
		return result, nil
	})

	// List alert rules
	server.RegisterTool(mcp.Tool{
		Name:        "grafana_list_alert_rules",
		Description: "List all Grafana alert rules",
		InputSchema: mcp.InputSchema{
			Type: "object",
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		rules, err := client.ListAlertRules(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list alert rules: %w", err)
		}
		return rules, nil
	})

	// Get health
	server.RegisterTool(mcp.Tool{
		Name:        "grafana_get_health",
		Description: "Check Grafana health status",
		InputSchema: mcp.InputSchema{
			Type: "object",
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		health, err := client.GetHealth(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get health: %w", err)
		}
		return health, nil
	})
}
