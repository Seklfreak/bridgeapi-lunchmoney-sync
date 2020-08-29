package lunchmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func (c *Client) GetAssets(ctx context.Context) ([]*Asset, error) {
	req, err := c.createRequest(ctx, http.MethodGet, "/v1/assets", nil)
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

	var assets assetsContainer
	err = json.Unmarshal(respData, &assets)
	if err != nil {
		return nil, err
	}

	return assets.Assets, nil
}

type assetsContainer struct {
	Assets []*Asset `json:"assets"`
}

type Asset struct {
	ID              int       `json:"id"`
	TypeName        string    `json:"type_name"`
	SubtypeName     string    `json:"subtype_name"`
	Name            string    `json:"name"`
	Balance         string    `json:"balance"`
	BalanceAsOf     time.Time `json:"balance_as_of"`
	Currency        string    `json:"currency"`
	InstitutionName string    `json:"institution_name"`
	CreatedAt       time.Time `json:"created_at"`
}
