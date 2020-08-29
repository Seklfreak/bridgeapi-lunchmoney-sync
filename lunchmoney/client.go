package lunchmoney

import (
	"bytes"
	"context"
	"net/http"
)

const baseURL = "https://dev.lunchmoney.app"

type Client struct {
	httpClient  *http.Client
	accessToken string
}

func NewClient(httpClient *http.Client, accessToken string) *Client {
	return &Client{
		httpClient:  httpClient,
		accessToken: accessToken,
	}
}

func (c *Client) createRequest(ctx context.Context, method string, endpoint string, body []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	return req, nil
}
