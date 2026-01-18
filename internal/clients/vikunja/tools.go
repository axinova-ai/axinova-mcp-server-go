package vikunja

import (
	"context"
	"fmt"
	"time"

	"github.com/axinova-ai/axinova-mcp-server-go/internal/mcp"
)

// RegisterTools registers all Vikunja tools with the MCP server
func RegisterTools(server *mcp.Server, client *Client) {
	// List projects
	server.RegisterTool(mcp.Tool{
		Name:        "vikunja_list_projects",
		Description: "List all Vikunja projects (lists)",
		InputSchema: mcp.InputSchema{
			Type: "object",
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		projects, err := client.ListProjects(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list projects: %w", err)
		}
		return projects, nil
	})

	// Get project
	server.RegisterTool(mcp.Tool{
		Name:        "vikunja_get_project",
		Description: "Get a specific Vikunja project by ID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"project_id": {
					Type:        "number",
					Description: "Project ID",
				},
			},
			Required: []string{"project_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		projectID, ok := args["project_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("project_id is required")
		}

		project, err := client.GetProject(ctx, int(projectID))
		if err != nil {
			return nil, fmt.Errorf("failed to get project: %w", err)
		}
		return project, nil
	})

	// Create project
	server.RegisterTool(mcp.Tool{
		Name:        "vikunja_create_project",
		Description: "Create a new Vikunja project",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"title": {
					Type:        "string",
					Description: "Project title",
				},
				"description": {
					Type:        "string",
					Description: "Project description (optional)",
				},
			},
			Required: []string{"title"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		title, ok := args["title"].(string)
		if !ok {
			return nil, fmt.Errorf("title is required")
		}

		description := ""
		if desc, ok := args["description"].(string); ok {
			description = desc
		}

		project, err := client.CreateProject(ctx, title, description)
		if err != nil {
			return nil, fmt.Errorf("failed to create project: %w", err)
		}
		return project, nil
	})

	// List tasks
	server.RegisterTool(mcp.Tool{
		Name:        "vikunja_list_tasks",
		Description: "List all tasks in a Vikunja project",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"project_id": {
					Type:        "number",
					Description: "Project ID",
				},
			},
			Required: []string{"project_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		projectID, ok := args["project_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("project_id is required")
		}

		tasks, err := client.ListTasks(ctx, int(projectID))
		if err != nil {
			return nil, fmt.Errorf("failed to list tasks: %w", err)
		}
		return tasks, nil
	})

	// Get task
	server.RegisterTool(mcp.Tool{
		Name:        "vikunja_get_task",
		Description: "Get a specific task by ID",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"project_id": {
					Type:        "number",
					Description: "Project ID",
				},
				"task_id": {
					Type:        "number",
					Description: "Task ID",
				},
			},
			Required: []string{"project_id", "task_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		projectID, ok := args["project_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("project_id is required")
		}

		taskID, ok := args["task_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("task_id is required")
		}

		task, err := client.GetTask(ctx, int(projectID), int(taskID))
		if err != nil {
			return nil, fmt.Errorf("failed to get task: %w", err)
		}
		return task, nil
	})

	// Create task
	server.RegisterTool(mcp.Tool{
		Name:        "vikunja_create_task",
		Description: "Create a new task in a Vikunja project",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"project_id": {
					Type:        "number",
					Description: "Project ID",
				},
				"title": {
					Type:        "string",
					Description: "Task title",
				},
				"description": {
					Type:        "string",
					Description: "Task description (optional)",
				},
				"priority": {
					Type:        "number",
					Description: "Task priority (0-5, default: 0)",
				},
				"due_date": {
					Type:        "string",
					Description: "Due date in RFC3339 format (optional)",
				},
			},
			Required: []string{"project_id", "title"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		projectID, ok := args["project_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("project_id is required")
		}

		title, ok := args["title"].(string)
		if !ok {
			return nil, fmt.Errorf("title is required")
		}

		req := CreateTaskRequest{
			Title: title,
		}

		if desc, ok := args["description"].(string); ok {
			req.Description = desc
		}

		if priority, ok := args["priority"].(float64); ok {
			req.Priority = int(priority)
		}

		if dueDate, ok := args["due_date"].(string); ok {
			t, err := time.Parse(time.RFC3339, dueDate)
			if err != nil {
				return nil, fmt.Errorf("invalid due_date format: %w", err)
			}
			req.DueDate = t
		}

		task, err := client.CreateTask(ctx, int(projectID), req)
		if err != nil {
			return nil, fmt.Errorf("failed to create task: %w", err)
		}
		return task, nil
	})

	// Update task
	server.RegisterTool(mcp.Tool{
		Name:        "vikunja_update_task",
		Description: "Update an existing Vikunja task",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"project_id": {
					Type:        "number",
					Description: "Project ID",
				},
				"task_id": {
					Type:        "number",
					Description: "Task ID",
				},
				"title": {
					Type:        "string",
					Description: "New task title (optional)",
				},
				"description": {
					Type:        "string",
					Description: "New task description (optional)",
				},
				"done": {
					Type:        "boolean",
					Description: "Mark task as done/undone (optional)",
				},
				"priority": {
					Type:        "number",
					Description: "New priority (0-5, optional)",
				},
			},
			Required: []string{"project_id", "task_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		projectID, ok := args["project_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("project_id is required")
		}

		taskID, ok := args["task_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("task_id is required")
		}

		req := UpdateTaskRequest{}

		if title, ok := args["title"].(string); ok {
			req.Title = title
		}

		if desc, ok := args["description"].(string); ok {
			req.Description = desc
		}

		if done, ok := args["done"].(bool); ok {
			req.Done = done
		}

		if priority, ok := args["priority"].(float64); ok {
			req.Priority = int(priority)
		}

		task, err := client.UpdateTask(ctx, int(projectID), int(taskID), req)
		if err != nil {
			return nil, fmt.Errorf("failed to update task: %w", err)
		}
		return task, nil
	})

	// Delete task
	server.RegisterTool(mcp.Tool{
		Name:        "vikunja_delete_task",
		Description: "Delete a Vikunja task",
		InputSchema: mcp.InputSchema{
			Type: "object",
			Properties: map[string]mcp.Property{
				"project_id": {
					Type:        "number",
					Description: "Project ID",
				},
				"task_id": {
					Type:        "number",
					Description: "Task ID",
				},
			},
			Required: []string{"project_id", "task_id"},
		},
	}, func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		projectID, ok := args["project_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("project_id is required")
		}

		taskID, ok := args["task_id"].(float64)
		if !ok {
			return nil, fmt.Errorf("task_id is required")
		}

		if err := client.DeleteTask(ctx, int(projectID), int(taskID)); err != nil {
			return nil, fmt.Errorf("failed to delete task: %w", err)
		}
		return fmt.Sprintf("Task %d deleted successfully", int(taskID)), nil
	})
}
