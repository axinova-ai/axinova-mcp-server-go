package grafana

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

// Client is a Grafana API client
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new Grafana client
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

// Dashboard represents a Grafana dashboard
type Dashboard struct {
	ID        int       `json:"id"`
	UID       string    `json:"uid"`
	Title     string    `json:"title"`
	Tags      []string  `json:"tags"`
	URL       string    `json:"url"`
	FolderID  int       `json:"folderId"`
	FolderUID string    `json:"folderUid"`
	IsStarred bool      `json:"isStarred"`
}

// DashboardMeta represents dashboard metadata
type DashboardMeta struct {
	IsStarred bool   `json:"isStarred"`
	Slug      string `json:"slug"`
	FolderID  int    `json:"folderId"`
	FolderUID string `json:"folderUid"`
}

// DashboardDetail represents full dashboard with panels
type DashboardDetail struct {
	Meta      DashboardMeta      `json:"meta"`
	Dashboard map[string]interface{} `json:"dashboard"`
}

// Datasource represents a Grafana datasource
type Datasource struct {
	ID       int    `json:"id"`
	UID      string `json:"uid"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	URL      string `json:"url"`
	IsDefault bool  `json:"isDefault"`
}

// AlertRule represents a Grafana alert rule
type AlertRule struct {
	ID          int    `json:"id"`
	UID         string `json:"uid"`
	Title       string `json:"title"`
	Condition   string `json:"condition"`
	Data        []AlertQuery `json:"data"`
	FolderUID   string `json:"folderUID"`
	RuleGroup   string `json:"ruleGroup"`
}

type AlertQuery struct {
	RefID         string                 `json:"refId"`
	QueryType     string                 `json:"queryType"`
	Model         map[string]interface{} `json:"model"`
	DatasourceUID string                 `json:"datasourceUid"`
}

// ListDashboards lists all dashboards
func (c *Client) ListDashboards(ctx context.Context) ([]Dashboard, error) {
	url := fmt.Sprintf("%s/api/search?type=dash-db", c.baseURL)

	var dashboards []Dashboard
	if err := c.doRequest(ctx, "GET", url, nil, &dashboards); err != nil {
		return nil, err
	}

	return dashboards, nil
}

// GetDashboard gets a dashboard by UID
func (c *Client) GetDashboard(ctx context.Context, uid string) (*DashboardDetail, error) {
	url := fmt.Sprintf("%s/api/dashboards/uid/%s", c.baseURL, uid)

	var detail DashboardDetail
	if err := c.doRequest(ctx, "GET", url, nil, &detail); err != nil {
		return nil, err
	}

	return &detail, nil
}

// CreateDashboard creates a new dashboard
func (c *Client) CreateDashboard(ctx context.Context, dashboard map[string]interface{}, folderUID string, overwrite bool) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/dashboards/db", c.baseURL)

	payload := map[string]interface{}{
		"dashboard": dashboard,
		"folderUid": folderUID,
		"overwrite": overwrite,
		"message":   "Created via MCP",
	}

	var result map[string]interface{}
	if err := c.doRequest(ctx, "POST", url, payload, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// DeleteDashboard deletes a dashboard by UID
func (c *Client) DeleteDashboard(ctx context.Context, uid string) error {
	url := fmt.Sprintf("%s/api/dashboards/uid/%s", c.baseURL, uid)
	return c.doRequest(ctx, "DELETE", url, nil, nil)
}

// ListDatasources lists all datasources
func (c *Client) ListDatasources(ctx context.Context) ([]Datasource, error) {
	url := fmt.Sprintf("%s/api/datasources", c.baseURL)

	var datasources []Datasource
	if err := c.doRequest(ctx, "GET", url, nil, &datasources); err != nil {
		return nil, err
	}

	return datasources, nil
}

// CreateDatasource creates a new datasource
func (c *Client) CreateDatasource(ctx context.Context, ds Datasource) (*Datasource, error) {
	url := fmt.Sprintf("%s/api/datasources", c.baseURL)

	var result Datasource
	if err := c.doRequest(ctx, "POST", url, ds, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// QueryDatasource queries a datasource (for Prometheus queries via Grafana)
func (c *Client) QueryDatasource(ctx context.Context, datasourceUID, query string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/ds/query", c.baseURL)

	payload := map[string]interface{}{
		"queries": []map[string]interface{}{
			{
				"refId":         "A",
				"datasourceUid": datasourceUID,
				"expr":          query,
				"format":        "time_series",
			},
		},
	}

	var result map[string]interface{}
	if err := c.doRequest(ctx, "POST", url, payload, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ListAlertRules lists all alert rules
func (c *Client) ListAlertRules(ctx context.Context) ([]AlertRule, error) {
	url := fmt.Sprintf("%s/api/v1/provisioning/alert-rules", c.baseURL)

	var rules []AlertRule
	if err := c.doRequest(ctx, "GET", url, nil, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

// GetHealth checks Grafana health
func (c *Client) GetHealth(ctx context.Context) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/health", c.baseURL)

	var health map[string]interface{}
	if err := c.doRequest(ctx, "GET", url, nil, &health); err != nil {
		return nil, err
	}

	return health, nil
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
