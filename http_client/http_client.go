package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"loro-tui/domain"
	"net/http"
	"time"
)

type Client struct {
	client *http.Client
	url    string
}

func NewClient(url string) (*Client, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}

	response, err := client.Get(url + "/health-check")
	if err != nil {
		return nil, err
	}
	if response.StatusCode != 200 {
		return nil, fmt.Errorf(response.Status)
	}

	return &Client{
		client: &client,
		url:    url,
	}, nil
}

func (c *Client) GetMessagesFrom(chatID int) {

}

func (c *Client) Login(payload RequestLogin) (*domain.UserInfo, error) {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	response, err := c.doRequest("POST", c.url+"/login", bytes, headers)
	if err != nil {
		return nil, err
	}
	userInfo := new(domain.UserInfo)
	err = json.Unmarshal(response, userInfo)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}

func (c *Client) doRequest(method, url string, body []byte, headers map[string]string) ([]byte, error) {
	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	for key, value := range headers {
		req.Header.Set(key, value)
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
