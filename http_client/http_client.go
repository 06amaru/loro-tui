package http_client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient(url string) (*Client, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	_, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: &client,
	}, nil
}

func (c *Client) Get(url string) ([]byte, error) {
	return c.doRequest(http.MethodGet, url, nil)
}

func (c *Client) Post(url string, body []byte) ([]byte, error) {
	return c.doRequest(http.MethodPost, url, body)
}

func (c *Client) doRequest(method, url string, body []byte) ([]byte, error) {
	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return respBody, nil
}
