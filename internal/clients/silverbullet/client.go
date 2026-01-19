package silverbullet

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client is a SilverBullet API client
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// NewClient creates a new SilverBullet client
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

// Page represents a SilverBullet page/note
type Page struct {
	Name         string    `json:"name"`
	LastModified time.Time `json:"lastModified"`
	Size         int       `json:"size"`
	ContentType  string    `json:"contentType"`
}

// SearchResult represents a search result
type SearchResult struct {
	Name    string `json:"name"`
	Text    string `json:"text"`
	Context string `json:"context"`
}

// ListPages lists all pages
func (c *Client) ListPages(ctx context.Context) ([]Page, error) {
	url := fmt.Sprintf("%s/index.json", c.baseURL)

	var pages []Page
	if err := c.doRequest(ctx, "GET", url, nil, &pages); err != nil {
		return nil, err
	}

	return pages, nil
}

// GetPage retrieves a page content
func (c *Client) GetPage(ctx context.Context, pageName string) (string, error) {
	url := fmt.Sprintf("%s/%s.md", c.baseURL, url.PathEscape(pageName))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// CreatePage creates a new page
func (c *Client) CreatePage(ctx context.Context, pageName, content string) error {
	url := fmt.Sprintf("%s/%s.md", c.baseURL, url.PathEscape(pageName))

	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBufferString(content))
	if err != nil {
		return err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Content-Type", "text/markdown")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// UpdatePage updates an existing page
func (c *Client) UpdatePage(ctx context.Context, pageName, content string) error {
	// Same as create - PUT is idempotent
	return c.CreatePage(ctx, pageName, content)
}

// DeletePage deletes a page
func (c *Client) DeletePage(ctx context.Context, pageName string) error {
	url := fmt.Sprintf("%s/%s.md", c.baseURL, url.PathEscape(pageName))

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// SearchPages searches pages by query
func (c *Client) SearchPages(ctx context.Context, query string) ([]SearchResult, error) {
	// Note: This is a basic implementation
	// SilverBullet's actual search API might differ
	url := fmt.Sprintf("%s/.client/search.json?query=%s", c.baseURL, url.QueryEscape(query))

	var results []SearchResult
	if err := c.doRequest(ctx, "GET", url, nil, &results); err != nil {
		// If search endpoint doesn't exist, fall back to listing all pages
		// and filtering by name
		pages, err := c.ListPages(ctx)
		if err != nil {
			return nil, err
		}

		results = []SearchResult{}
		for _, page := range pages {
			results = append(results, SearchResult{
				Name: page.Name,
				Text: "",
			})
		}
	}

	return results, nil
}

// doRequest performs an HTTP request with JSON response
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

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
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
