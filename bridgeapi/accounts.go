package bridgeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (c *Client) fetchAccounts(ctx context.Context, endpoint string) (*AccountsContainer, error) {
	req, err := c.createRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var accounts AccountsContainer
	err = json.Unmarshal(respData, &accounts)
	if err != nil {
		return nil, err
	}

	return &accounts, nil
}

func (c *Client) FetchAccounts(ctx context.Context) ([]*Account, error) {
	var accounts []*Account

	var endpoint string
	for {
		if endpoint == "" {
			endpoint = "/v2/accounts"
		}

		result, err := c.fetchAccounts(ctx, endpoint)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, result.Accounts...)

		endpoint = result.Pagination.NextURI

		if len(result.Accounts) <= 0 || endpoint == "" {
			break
		}
	}

	return accounts, nil
}

type AccountsContainer struct {
	Accounts   []*Account `json:"resources"`
	Pagination Pagination `json:"pagination"`
}

type Account struct {
	ID                    int         `json:"id"`
	ResourceURI           string      `json:"resource_uri"`
	ResourceType          string      `json:"resource_type"`
	Name                  string      `json:"name"`
	Balance               float64     `json:"balance"`
	Status                int         `json:"status"`
	StatusCodeInfo        string      `json:"status_code_info"`
	StatusCodeDescription interface{} `json:"status_code_description"`
	UpdatedAt             time.Time   `json:"updated_at"`
	Type                  string      `json:"type"`
	CurrencyCode          string      `json:"currency_code"`
	Item                  Item        `json:"item"`
	Bank                  Bank        `json:"bank"`
	LoanDetails           interface{} `json:"loan_details"`
	SavingsDetails        interface{} `json:"savings_details"`
	IsPro                 bool        `json:"is_pro"`
	LastRefresh           time.Time   `json:"last_refresh"`
	Iban                  interface{} `json:"iban"`
	IbanFillManually      bool        `json:"iban_fill_manually"`
}

type Item struct {
	ID           int    `json:"id"`
	ResourceURI  string `json:"resource_uri"`
	ResourceType string `json:"resource_type"`
}

type Bank struct {
	ID           int    `json:"id"`
	ResourceURI  string `json:"resource_uri"`
	ResourceType string `json:"resource_type"`
}
