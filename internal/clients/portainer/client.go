package portainer

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client is a Portainer API client
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new Portainer client
func NewClient(baseURL, token string, timeout time.Duration, skipTLSVerify bool) *Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipTLSVerify,
		},
	}

	return &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout:   timeout,
			Transport: transport,
		},
	}
}

// Container represents a Docker container
type Container struct {
	Id      string            `json:"Id"`
	Names   []string          `json:"Names"`
	Image   string            `json:"Image"`
	State   string            `json:"State"`
	Status  string            `json:"Status"`
	Labels  map[string]string `json:"Labels"`
}

// Stack represents a Docker Compose stack
type Stack struct {
	Id             int      `json:"Id"`
	Name           string   `json:"Name"`
	Type           int      `json:"Type"`
	EndpointId     int      `json:"EndpointId"`
	SwarmId        string   `json:"SwarmId,omitempty"`
	EntryPoint     string   `json:"EntryPoint"`
	Env            []StackEnv `json:"Env"`
	Status         int      `json:"Status"`
}

type StackEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ContainerStats represents container resource statistics
type ContainerStats struct {
	CPUPercent    float64 `json:"cpu_percent"`
	MemoryUsage   uint64  `json:"memory_usage"`
	MemoryLimit   uint64  `json:"memory_limit"`
	MemoryPercent float64 `json:"memory_percent"`
	NetworkRx     uint64  `json:"network_rx"`
	NetworkTx     uint64  `json:"network_tx"`
}

// ListContainers lists all containers
func (c *Client) ListContainers(ctx context.Context, endpointID int) ([]Container, error) {
	url := fmt.Sprintf("%s/api/endpoints/%d/docker/containers/json?all=1", c.baseURL, endpointID)

	var containers []Container
	if err := c.doRequest(ctx, "GET", url, nil, &containers); err != nil {
		return nil, err
	}

	return containers, nil
}

// StartContainer starts a container
func (c *Client) StartContainer(ctx context.Context, endpointID int, containerID string) error {
	url := fmt.Sprintf("%s/api/endpoints/%d/docker/containers/%s/start", c.baseURL, endpointID, containerID)
	return c.doRequest(ctx, "POST", url, nil, nil)
}

// StopContainer stops a container
func (c *Client) StopContainer(ctx context.Context, endpointID int, containerID string) error {
	url := fmt.Sprintf("%s/api/endpoints/%d/docker/containers/%s/stop", c.baseURL, endpointID, containerID)
	return c.doRequest(ctx, "POST", url, nil, nil)
}

// RestartContainer restarts a container
func (c *Client) RestartContainer(ctx context.Context, endpointID int, containerID string) error {
	url := fmt.Sprintf("%s/api/endpoints/%d/docker/containers/%s/restart", c.baseURL, endpointID, containerID)
	return c.doRequest(ctx, "POST", url, nil, nil)
}

// GetContainerLogs retrieves container logs
func (c *Client) GetContainerLogs(ctx context.Context, endpointID int, containerID string, tail int) (string, error) {
	url := fmt.Sprintf("%s/api/endpoints/%d/docker/containers/%s/logs?stdout=1&stderr=1&tail=%d",
		c.baseURL, endpointID, containerID, tail)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("X-API-Key", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	logs, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(logs), nil
}

// ListStacks lists all stacks
func (c *Client) ListStacks(ctx context.Context) ([]Stack, error) {
	url := fmt.Sprintf("%s/api/stacks", c.baseURL)

	var stacks []Stack
	if err := c.doRequest(ctx, "GET", url, nil, &stacks); err != nil {
		return nil, err
	}

	return stacks, nil
}

// GetStack gets a stack by ID
func (c *Client) GetStack(ctx context.Context, stackID int) (*Stack, error) {
	url := fmt.Sprintf("%s/api/stacks/%d", c.baseURL, stackID)

	var stack Stack
	if err := c.doRequest(ctx, "GET", url, nil, &stack); err != nil {
		return nil, err
	}

	return &stack, nil
}

// InspectContainer gets detailed container info
func (c *Client) InspectContainer(ctx context.Context, endpointID int, containerID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/endpoints/%d/docker/containers/%s/json", c.baseURL, endpointID, containerID)

	var info map[string]interface{}
	if err := c.doRequest(ctx, "GET", url, nil, &info); err != nil {
		return nil, err
	}

	return info, nil
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

	req.Header.Set("X-API-Key", c.token)
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

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	return nil
}
