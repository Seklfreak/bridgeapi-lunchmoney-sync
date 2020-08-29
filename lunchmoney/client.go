package lunchmoney

import (
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
