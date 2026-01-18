package vikunja

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a Vikunja API client
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new Vikunja client
func NewClient(baseURL, token string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// Project represents a Vikunja project (list)
type Project struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
}

// Task represents a Vikunja task
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Done        bool      `json:"done"`
	ProjectID   int       `json:"project_id"`
	Priority    int       `json:"priority"`
	DueDate     time.Time `json:"due_date,omitempty"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	CreatedBy   User      `json:"created_by"`
	Labels      []Label   `json:"labels,omitempty"`
}

// User represents a Vikunja user
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// Label represents a task label
type Label struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Color string `json:"hex_color"`
}

// CreateTaskRequest represents a task creation request
type CreateTaskRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Done        bool      `json:"done,omitempty"`
	Priority    int       `json:"priority,omitempty"`
	DueDate     time.Time `json:"due_date,omitempty"`
}

// UpdateTaskRequest represents a task update request
type UpdateTaskRequest struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Done        bool      `json:"done,omitempty"`
	Priority    int       `json:"priority,omitempty"`
	DueDate     time.Time `json:"due_date,omitempty"`
}

// ListProjects lists all projects
func (c *Client) ListProjects(ctx context.Context) ([]Project, error) {
	url := fmt.Sprintf("%s/api/v1/projects", c.baseURL)

	var projects []Project
	if err := c.doRequest(ctx, "GET", url, nil, &projects); err != nil {
		return nil, err
	}

	return projects, nil
}

// GetProject gets a project by ID
func (c *Client) GetProject(ctx context.Context, projectID int) (*Project, error) {
	url := fmt.Sprintf("%s/api/v1/projects/%d", c.baseURL, projectID)

	var project Project
	if err := c.doRequest(ctx, "GET", url, nil, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

// CreateProject creates a new project
func (c *Client) CreateProject(ctx context.Context, title, description string) (*Project, error) {
	url := fmt.Sprintf("%s/api/v1/projects", c.baseURL)

	payload := map[string]interface{}{
		"title":       title,
		"description": description,
	}

	var project Project
	if err := c.doRequest(ctx, "POST", url, payload, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

// ListTasks lists all tasks in a project
func (c *Client) ListTasks(ctx context.Context, projectID int) ([]Task, error) {
	url := fmt.Sprintf("%s/api/v1/projects/%d/tasks", c.baseURL, projectID)

	var tasks []Task
	if err := c.doRequest(ctx, "GET", url, nil, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

// GetTask gets a task by ID
func (c *Client) GetTask(ctx context.Context, projectID, taskID int) (*Task, error) {
	url := fmt.Sprintf("%s/api/v1/projects/%d/tasks/%d", c.baseURL, projectID, taskID)

	var task Task
	if err := c.doRequest(ctx, "GET", url, nil, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// CreateTask creates a new task
func (c *Client) CreateTask(ctx context.Context, projectID int, req CreateTaskRequest) (*Task, error) {
	url := fmt.Sprintf("%s/api/v1/projects/%d/tasks", c.baseURL, projectID)

	var task Task
	if err := c.doRequest(ctx, "PUT", url, req, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// UpdateTask updates a task
func (c *Client) UpdateTask(ctx context.Context, projectID, taskID int, req UpdateTaskRequest) (*Task, error) {
	url := fmt.Sprintf("%s/api/v1/projects/%d/tasks/%d", c.baseURL, projectID, taskID)

	var task Task
	if err := c.doRequest(ctx, "POST", url, req, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// DeleteTask deletes a task
func (c *Client) DeleteTask(ctx context.Context, projectID, taskID int) error {
	url := fmt.Sprintf("%s/api/v1/projects/%d/tasks/%d", c.baseURL, projectID, taskID)
	return c.doRequest(ctx, "DELETE", url, nil, nil)
}

// doRequest performs an HTTP request
func (c *Client) doRequest(ctx context.Context, method, url string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil && resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}
