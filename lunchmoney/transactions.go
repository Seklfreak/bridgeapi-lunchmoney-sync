package lunchmoney

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Transaction struct {
	Date        string   `json:"date"`
	Amount      float64  `json:"amount"`
	CategoryID  int      `json:"category_id,omitempty"`
	Payee       string   `json:"payee,omitempty"`
	Currency    string   `json:"currency,omitempty"`
	AssetID     int      `json:"asset_id,omitempty"`
	RecurringID int      `json:"recurring_id,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	Status      string   `json:"status,omitempty"`
	ExternalID  string   `json:"external_id,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (c *Client) InsertTransactions(ctx context.Context, trx []*Transaction) (int, error) {
	var request struct {
		Transactions      []*Transaction `json:"transactions"`
		ApplyRules        bool           `json:"apply_rules"`
		CheckForRecurring bool           `json:"check_for_recurring"`
		DebitAsNegative   bool           `json:"debit_as_negative"`
	}
	request.Transactions = trx
	request.ApplyRules = true
	request.CheckForRecurring = true
	request.DebitAsNegative = true

	reqData, err := json.Marshal(request)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v1/transactions", bytes.NewReader(reqData))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Error []string `json:"error"`
		IDs   []int    `json:"ids"`
	}
	err = json.Unmarshal(respData, &result)
	if err != nil {
		return 0, err
	}

	if len(result.Error) > 0 {
		return len(result.IDs), fmt.Errorf("received %d errors: %q", len(result.Error), strings.Join(result.Error, "; "))
	}

	return len(result.IDs), nil
}
