package lunchmoney

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Transaction struct {
	Date        string   `json:"date"`
	Amount      float64  `json:"amount"`
	CategoryID  int      `json:"category_id,omitempty"`
	Payee       string   `json:"payee,omitempty"`
	Currency    string   `json:"currency,omitempty"`
	AssetID     int      `json:"asset,omitempty"`
	RecurringID int      `json:"recurring_id,omitempty"`
	Notes       string   `json:"notes,omitempty"`
	Status      string   `json:"status,omitempty"`
	ExternalID  string   `json:"external_id,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

func (c *Client) InsertTransactions(ctx context.Context, trx []*Transaction) error {
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
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/v1/transactions", bytes.NewReader(reqData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_ = respData

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
