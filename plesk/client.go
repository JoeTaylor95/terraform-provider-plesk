package plesk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

type Client struct {
	Host   string
	Port   string
	Token  string
	Client *http.Client
}

func (c *Client) Get(ctx context.Context, path string) ([]byte, diag.Diagnostics) {
	url := fmt.Sprintf("https://%s:%s%s", c.Host, c.Port, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to create GET request: %s", err),
			},
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("GET request failed: %s", err),
			},
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to read GET response body: %s", err),
			},
		}
	}

	if resp.StatusCode >= 400 {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("GET request returned HTTP %d: %s", resp.StatusCode, string(body)),
			},
		}
	}

	return body, nil
}

func (c *Client) Post(ctx context.Context, path string, data interface{}) ([]byte, diag.Diagnostics) {
	url := fmt.Sprintf("https://%s:%s%s", c.Host, c.Port, path)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to marshal POST data: %s", err),
			},
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to create POST request: %s", err),
			},
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("POST request failed: %s", err),
			},
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to read POST response body: %s", err),
			},
		}
	}

	if resp.StatusCode >= 400 {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("POST request returned HTTP %d: %s", resp.StatusCode, string(body)),
			},
		}
	}

	return body, nil
}

func (c *Client) Put(ctx context.Context, path string, data interface{}) ([]byte, diag.Diagnostics) {
	url := fmt.Sprintf("https://%s:%s%s", c.Host, c.Port, path)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to marshal PUT data: %s", err),
			},
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to create PUT request: %s", err),
			},
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("PUT request failed: %s", err),
			},
		}
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to read PUT response body: %s", err),
			},
		}
	}

	if resp.StatusCode >= 400 {
		return nil, diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("PUT request returned HTTP %d: %s", resp.StatusCode, string(body)),
			},
		}
	}

	return body, nil
}

func (c *Client) Delete(ctx context.Context, path string) diag.Diagnostics {
	url := fmt.Sprintf("https://%s:%s%s", c.Host, c.Port, path)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to create DELETE request: %s", err),
			},
		}
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("DELETE request failed: %s", err),
			},
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := ioutil.ReadAll(resp.Body)
		return diag.Diagnostics{
			{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("DELETE request returned HTTP %d: %s", resp.StatusCode, string(body)),
			},
		}
	}

	return nil
}
