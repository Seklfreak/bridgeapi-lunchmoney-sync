package lunchmoney

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

func (c *Client) UpdateAsset(ctx context.Context, assetID int, asset *Asset) error {
	reqData, err := json.Marshal(asset)
	if err != nil {
		return err
	}

	req, err := c.createRequest(ctx, http.MethodPut, "/v1/assets/"+strconv.Itoa(assetID), reqData)
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
		Error []string `json:"error"`
	}
	err = json.Unmarshal(respData, &result)
	if err != nil {
		return err
	}

	if len(result.Error) > 0 {
		return fmt.Errorf("received %d errors: %q", len(result.Error), strings.Join(result.Error, "; "))
	}

	return nil
}

type assetsContainer struct {
	Assets []*Asset `json:"assets"`
}

type Asset struct {
	ID              int        `json:"id,omitempty"`
	TypeName        string     `json:"type_name,omitempty"`
	SubtypeName     string     `json:"subtype_name,omitempty"`
	Name            string     `json:"name,omitempty"`
	Balance         string     `json:"balance,omitempty"`
	BalanceAsOf     *time.Time `json:"balance_as_of,omitempty"`
	Currency        string     `json:"currency,omitempty"`
	InstitutionName string     `json:"institution_name,omitempty"`
	CreatedAt       *time.Time `json:"created_at,omitempty"`
}
