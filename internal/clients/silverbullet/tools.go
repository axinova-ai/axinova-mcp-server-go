package silverbullet

import (
	"context"
	"fmt"

	"github.com/axinova-ai/axinova-mcp-server-go/internal/mcp"
)

// RegisterTools registers all SilverBullet tools with the MCP server
func RegisterTools(server *mcp.Server, client *Client) {
	// List pages
	server.RegisterTool(mcp.Tool{
		Name:        "silverbullet_list_pages",
		Description: "List all SilverBullet pages/notes",
		InputSchema: mcp.InputSchema{
			Type: "object",
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		pages, err := client.ListPages(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list pages: %w", err)
		}
		return pages, nil
	})

	// Get page
	server.RegisterTool(mcp.Tool{
		Name:        "silverbullet_get_page",
		Description: "Get content of a specific SilverBullet page",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"page_name": {
					Type:        "string",
					Description: "Page name (without .md extension)",
				},
			},
			Required: []string{"page_name"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		pageName, ok := args["page_name"].(string)
		if !ok {
			return nil, fmt.Errorf("page_name is required")
		}

		content, err := client.GetPage(ctx, pageName)
		if err != nil {
			return nil, fmt.Errorf("failed to get page: %w", err)
		}
		return content, nil
	})

	// Create page
	server.RegisterTool(mcp.Tool{
		Name:        "silverbullet_create_page",
		Description: "Create a new SilverBullet page/note",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"page_name": {
					Type:        "string",
					Description: "Page name (without .md extension)",
				},
				"content": {
					Type:        "string",
					Description: "Page content in Markdown format",
				},
			},
			Required: []string{"page_name", "content"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		pageName, ok := args["page_name"].(string)
		if !ok {
			return nil, fmt.Errorf("page_name is required")
		}

		content, ok := args["content"].(string)
		if !ok {
			return nil, fmt.Errorf("content is required")
		}

		if err := client.CreatePage(ctx, pageName, content); err != nil {
			return nil, fmt.Errorf("failed to create page: %w", err)
		}
		return fmt.Sprintf("Page '%s' created successfully", pageName), nil
	})

	// Update page
	server.RegisterTool(mcp.Tool{
		Name:        "silverbullet_update_page",
		Description: "Update an existing SilverBullet page",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"page_name": {
					Type:        "string",
					Description: "Page name (without .md extension)",
				},
				"content": {
					Type:        "string",
					Description: "New page content in Markdown format",
				},
			},
			Required: []string{"page_name", "content"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		pageName, ok := args["page_name"].(string)
		if !ok {
			return nil, fmt.Errorf("page_name is required")
		}

		content, ok := args["content"].(string)
		if !ok {
			return nil, fmt.Errorf("content is required")
		}

		if err := client.UpdatePage(ctx, pageName, content); err != nil {
			return nil, fmt.Errorf("failed to update page: %w", err)
		}
		return fmt.Sprintf("Page '%s' updated successfully", pageName), nil
	})

	// Delete page
	server.RegisterTool(mcp.Tool{
		Name:        "silverbullet_delete_page",
		Description: "Delete a SilverBullet page",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"page_name": {
					Type:        "string",
					Description: "Page name (without .md extension)",
				},
			},
			Required: []string{"page_name"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		pageName, ok := args["page_name"].(string)
		if !ok {
			return nil, fmt.Errorf("page_name is required")
		}

		if err := client.DeletePage(ctx, pageName); err != nil {
			return nil, fmt.Errorf("failed to delete page: %w", err)
		}
		return fmt.Sprintf("Page '%s' deleted successfully", pageName), nil
	})

	// Search pages
	server.RegisterTool(mcp.Tool{
		Name:        "silverbullet_search_pages",
		Description: "Search SilverBullet pages by query",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"query": {
					Type:        "string",
					Description: "Search query",
				},
			},
			Required: []string{"query"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		query, ok := args["query"].(string)
		if !ok {
			return nil, fmt.Errorf("query is required")
		}

		results, err := client.SearchPages(ctx, query)
		if err != nil {
			return nil, fmt.Errorf("failed to search pages: %w", err)
		}
		return results, nil
	})
}
