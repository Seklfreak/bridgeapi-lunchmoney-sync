package bridgeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func (c *Client) fetchTransactions(ctx context.Context, endpoint string) (*TransactionsContainer, error) {
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

	var transactions TransactionsContainer
	err = json.Unmarshal(respData, &transactions)
	if err != nil {
		return nil, err
	}

	return &transactions, nil
}

func (c *Client) FetchTransactionsUpdated(ctx context.Context, since time.Time) ([]*Transaction, error) {
	var transactions []*Transaction

	var endpoint string
	for {
		if endpoint == "" {
			vars := url.Values{}
			vars.Set("since", since.UTC().Format(time.RFC3339))

			endpoint = "/v2/transactions/updated?" + vars.Encode()
		}

		result, err := c.fetchTransactions(ctx, endpoint)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, result.Transactions...)

		endpoint = result.Pagination.NextURI

		if len(result.Transactions) <= 0 || endpoint == "" {
			break
		}
	}

	return transactions, nil
}

type TransactionsContainer struct {
	Transactions []*Transaction `json:"resources"`
	Pagination   Pagination     `json:"pagination"`
}

type Transaction struct {
	ID             int64     `json:"id"`
	ResourceURI    string    `json:"resource_uri"`
	ResourceType   string    `json:"resource_type"`
	Description    string    `json:"description"`
	RawDescription string    `json:"raw_description"`
	Amount         float64   `json:"amount"`
	Date           string    `json:"date"`
	UpdatedAt      time.Time `json:"updated_at"`
	CurrencyCode   string    `json:"currency_code"`
	IsDeleted      bool      `json:"is_deleted"`
	Category       Category  `json:"category"`
	Account        Account   `json:"account"`
	IsFuture       bool      `json:"is_future"`
}

type Category struct {
	ID           int    `json:"id"`
	ResourceURI  string `json:"resource_uri"`
	ResourceType string `json:"resource_type"`
}

type Account struct {
	ID           int    `json:"id"`
	ResourceURI  string `json:"resource_uri"`
	ResourceType string `json:"resource_type"`
}

type Pagination struct {
	PreviousURI string `json:"previous_uri"`
	NextURI     string `json:"next_uri"`
}
