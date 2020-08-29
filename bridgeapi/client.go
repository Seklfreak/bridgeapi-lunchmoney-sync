package bridgeapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	baseURL    = "https://sync.bankin.com"
	apiVersion = "2019-02-18"
)

type Client struct {
	httpClient *http.Client
	auth       *Auth
}

type Auth struct {
	ClientID     string
	ClientSecret string
	Email        string
	Password     string

	accessToken string
}

func NewClient(httpClient *http.Client, auth *Auth) (*Client, error) {
	client := &Client{
		httpClient: httpClient,
		auth:       auth,
	}
	if client.auth.accessToken == "" {
		err := client.login()
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (c *Client) createRequest(method string, endpoint string, body []byte) (*http.Request, error) {
	req, err := http.NewRequest(method, baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Bankin-Version", apiVersion)
	req.Header.Set("Client-Id", c.auth.ClientID)
	req.Header.Set("Client-Secret", c.auth.ClientSecret)
	if c.auth.accessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.auth.accessToken)
	}

	return req, nil
}

func (c *Client) login() error {
	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	request.Email = c.auth.Email
	request.Password = c.auth.Password

	reqData, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := c.createRequest(http.MethodPost, "/v2/authenticate", reqData)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	err = json.Unmarshal(respData, &result)
	if err != nil {
		return err
	}
	if result.AccessToken == "" {
		return errors.New("login failed, access token is empty")
	}

	c.auth.accessToken = result.AccessToken
	return nil
}
