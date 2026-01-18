package prometheus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is a Prometheus API client
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new Prometheus client
func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// QueryResult represents a Prometheus query result
type QueryResult struct {
	Status string    `json:"status"`
	Data   QueryData `json:"data"`
}

type QueryData struct {
	ResultType string   `json:"resultType"`
	Result     []Metric `json:"result"`
}

type Metric struct {
	Metric map[string]string `json:"metric"`
	Value  []interface{}     `json:"value,omitempty"`
	Values [][]interface{}   `json:"values,omitempty"`
}

// LabelValues represents label values response
type LabelValues struct {
	Status string   `json:"status"`
	Data   []string `json:"data"`
}

// TargetsResult represents targets response
type TargetsResult struct {
	Status string       `json:"status"`
	Data   TargetsData  `json:"data"`
}

type TargetsData struct {
	ActiveTargets []Target `json:"activeTargets"`
	DroppedTargets []Target `json:"droppedTargets"`
}

type Target struct {
	DiscoveredLabels map[string]string `json:"discoveredLabels"`
	Labels           map[string]string `json:"labels"`
	ScrapePool       string            `json:"scrapePool"`
	ScrapeURL        string            `json:"scrapeUrl"`
	LastError        string            `json:"lastError"`
	LastScrape       time.Time         `json:"lastScrape"`
	Health           string            `json:"health"`
}

// Query executes an instant query
func (c *Client) Query(ctx context.Context, query string, timestamp *time.Time) (*QueryResult, error) {
	params := url.Values{}
	params.Add("query", query)
	if timestamp != nil {
		params.Add("time", fmt.Sprintf("%d", timestamp.Unix()))
	}

	url := fmt.Sprintf("%s/api/v1/query?%s", c.baseURL, params.Encode())

	var result QueryResult
	if err := c.doRequest(ctx, "GET", url, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// QueryRange executes a range query
func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time, step string) (*QueryResult, error) {
	params := url.Values{}
	params.Add("query", query)
	params.Add("start", fmt.Sprintf("%d", start.Unix()))
	params.Add("end", fmt.Sprintf("%d", end.Unix()))
	params.Add("step", step)

	url := fmt.Sprintf("%s/api/v1/query_range?%s", c.baseURL, params.Encode())

	var result QueryResult
	if err := c.doRequest(ctx, "GET", url, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// LabelNames returns all label names
func (c *Client) LabelNames(ctx context.Context) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/labels", c.baseURL)

	var result LabelValues
	if err := c.doRequest(ctx, "GET", url, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// LabelValues returns all label values for a given label name
func (c *Client) LabelValues(ctx context.Context, label string) ([]string, error) {
	url := fmt.Sprintf("%s/api/v1/label/%s/values", c.baseURL, label)

	var result LabelValues
	if err := c.doRequest(ctx, "GET", url, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// Series finds series by label matchers
func (c *Client) Series(ctx context.Context, matches []string, start, end time.Time) ([]map[string]string, error) {
	params := url.Values{}
	for _, match := range matches {
		params.Add("match[]", match)
	}
	params.Add("start", fmt.Sprintf("%d", start.Unix()))
	params.Add("end", fmt.Sprintf("%d", end.Unix()))

	url := fmt.Sprintf("%s/api/v1/series?%s", c.baseURL, params.Encode())

	var result struct {
		Status string              `json:"status"`
		Data   []map[string]string `json:"data"`
	}
	if err := c.doRequest(ctx, "GET", url, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// Targets returns all active and dropped targets
func (c *Client) Targets(ctx context.Context) (*TargetsResult, error) {
	url := fmt.Sprintf("%s/api/v1/targets", c.baseURL)

	var result TargetsResult
	if err := c.doRequest(ctx, "GET", url, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Metadata returns metric metadata
func (c *Client) Metadata(ctx context.Context, metric string) (map[string]interface{}, error) {
	params := url.Values{}
	if metric != "" {
		params.Add("metric", metric)
	}

	url := fmt.Sprintf("%s/api/v1/metadata?%s", c.baseURL, params.Encode())

	var result struct {
		Status string                 `json:"status"`
		Data   map[string]interface{} `json:"data"`
	}
	if err := c.doRequest(ctx, "GET", url, &result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

// doRequest performs an HTTP request
func (c *Client) doRequest(ctx context.Context, method, url string, result interface{}) error {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return err
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
