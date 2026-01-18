package portainer

import (
	"context"
	"fmt"
	"strconv"

	"github.com/axinova-ai/axinova-mcp-server-go/internal/mcp"
)

// RegisterTools registers all Portainer tools with the MCP server
func RegisterTools(server *mcp.Server, client *Client) {
	// List containers
	server.RegisterTool(mcp.Tool{
		Name:        "portainer_list_containers",
		Description: "List all Docker containers in a Portainer environment",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"endpoint_id": {
					Type:        "number",
					Description: "Portainer endpoint ID (default: 1 for local)",
				},
			},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		endpointID := 1
		if id, ok := args["endpoint_id"].(float64); ok {
			endpointID = int(id)
		}

		containers, err := client.ListContainers(ctx, endpointID)
		if err != nil {
			return nil, fmt.Errorf("failed to list containers: %w", err)
		}

		return containers, nil
	})

	// Start container
	server.RegisterTool(mcp.Tool{
		Name:        "portainer_start_container",
		Description: "Start a Docker container",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"endpoint_id": {
					Type:        "number",
					Description: "Portainer endpoint ID (default: 1)",
				},
				"container_id": {
					Type:        "string",
					Description: "Container ID or name",
				},
			},
			Required: []string{"container_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		endpointID := 1
		if id, ok := args["endpoint_id"].(float64); ok {
			endpointID = int(id)
		}

		containerID, ok := args["container_id"].(string)
		if !ok {
			return nil, fmt.Errorf("container_id is required")
		}

		if err := client.StartContainer(ctx, endpointID, containerID); err != nil {
			return nil, fmt.Errorf("failed to start container: %w", err)
		}

		return fmt.Sprintf("Container %s started successfully", containerID), nil
	})

	// Stop container
	server.RegisterTool(mcp.Tool{
		Name:        "portainer_stop_container",
		Description: "Stop a Docker container",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"endpoint_id": {
					Type:        "number",
					Description: "Portainer endpoint ID (default: 1)",
				},
				"container_id": {
					Type:        "string",
					Description: "Container ID or name",
				},
			},
			Required: []string{"container_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		endpointID := 1
		if id, ok := args["endpoint_id"].(float64); ok {
			endpointID = int(id)
		}

		containerID, ok := args["container_id"].(string)
		if !ok {
			return nil, fmt.Errorf("container_id is required")
		}

		if err := client.StopContainer(ctx, endpointID, containerID); err != nil {
			return nil, fmt.Errorf("failed to stop container: %w", err)
		}

		return fmt.Sprintf("Container %s stopped successfully", containerID), nil
	})

	// Restart container
	server.RegisterTool(mcp.Tool{
		Name:        "portainer_restart_container",
		Description: "Restart a Docker container",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"endpoint_id": {
					Type:        "number",
					Description: "Portainer endpoint ID (default: 1)",
				},
				"container_id": {
					Type:        "string",
					Description: "Container ID or name",
				},
			},
			Required: []string{"container_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		endpointID := 1
		if id, ok := args["endpoint_id"].(float64); ok {
			endpointID = int(id)
		}

		containerID, ok := args["container_id"].(string)
		if !ok {
			return nil, fmt.Errorf("container_id is required")
		}

		if err := client.RestartContainer(ctx, endpointID, containerID); err != nil {
			return nil, fmt.Errorf("failed to restart container: %w", err)
		}

		return fmt.Sprintf("Container %s restarted successfully", containerID), nil
	})

	// Get container logs
	server.RegisterTool(mcp.Tool{
		Name:        "portainer_get_container_logs",
		Description: "Retrieve logs from a Docker container",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"endpoint_id": {
					Type:        "number",
					Description: "Portainer endpoint ID (default: 1)",
				},
				"container_id": {
					Type:        "string",
					Description: "Container ID or name",
				},
				"tail": {
					Type:        "number",
					Description: "Number of log lines to retrieve (default: 100)",
				},
			},
			Required: []string{"container_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		endpointID := 1
		if id, ok := args["endpoint_id"].(float64); ok {
			endpointID = int(id)
		}

		containerID, ok := args["container_id"].(string)
		if !ok {
			return nil, fmt.Errorf("container_id is required")
		}

		tail := 100
		if t, ok := args["tail"].(float64); ok {
			tail = int(t)
		}

		logs, err := client.GetContainerLogs(ctx, endpointID, containerID, tail)
		if err != nil {
			return nil, fmt.Errorf("failed to get logs: %w", err)
		}

		return logs, nil
	})

	// List stacks
	server.RegisterTool(mcp.Tool{
		Name:        "portainer_list_stacks",
		Description: "List all Docker Compose stacks",
		InputSchema: mcp.InputSchema{
			Type: "object",
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		stacks, err := client.ListStacks(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list stacks: %w", err)
		}

		return stacks, nil
	})

	// Get stack details
	server.RegisterTool(mcp.Tool{
		Name:        "portainer_get_stack",
		Description: "Get details of a specific Docker Compose stack",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"stack_id": {
					Type:        "number",
					Description: "Stack ID",
				},
			},
			Required: []string{"stack_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		var stackID int
		switch v := args["stack_id"].(type) {
		case float64:
			stackID = int(v)
		case string:
			id, err := strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("invalid stack_id: %w", err)
			}
			stackID = id
		default:
			return nil, fmt.Errorf("stack_id is required")
		}

		stack, err := client.GetStack(ctx, stackID)
		if err != nil {
			return nil, fmt.Errorf("failed to get stack: %w", err)
		}

		return stack, nil
	})

	// Inspect container
	server.RegisterTool(mcp.Tool{
		Name:        "portainer_inspect_container",
		Description: "Get detailed information about a container",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"endpoint_id": {
					Type:        "number",
					Description: "Portainer endpoint ID (default: 1)",
				},
				"container_id": {
					Type:        "string",
					Description: "Container ID or name",
				},
			},
			Required: []string{"container_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		endpointID := 1
		if id, ok := args["endpoint_id"].(float64); ok {
			endpointID = int(id)
		}

		containerID, ok := args["container_id"].(string)
		if !ok {
			return nil, fmt.Errorf("container_id is required")
		}

		info, err := client.InspectContainer(ctx, endpointID, containerID)
		if err != nil {
			return nil, fmt.Errorf("failed to inspect container: %w", err)
		}

		return info, nil
	})
}
