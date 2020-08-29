package bridgeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (c *Client) FetchBank(ctx context.Context, id int) (*Bank, error) {
	req, err := c.createRequest(ctx, http.MethodGet, "/v2/banks/"+strconv.Itoa(id), nil)
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

	var bank Bank
	err = json.Unmarshal(respData, &bank)
	if err != nil {
		return nil, err
	}

	return &bank, nil
}

type Bank struct {
	ID                 int         `json:"id"`
	ResourceURI        string      `json:"resource_uri"`
	ResourceType       string      `json:"resource_type"`
	Name               string      `json:"name"`
	CountryCode        string      `json:"country_code"`
	AutomaticRefresh   bool        `json:"automatic_refresh"`
	PrimaryColor       interface{} `json:"primary_color"`
	SecondaryColor     interface{} `json:"secondary_color"`
	LogoURL            string      `json:"logo_url"`
	DeeplinkIos        interface{} `json:"deeplink_ios"`
	DeeplinkAndroid    interface{} `json:"deeplink_android"`
	TransferEnabled    bool        `json:"transfer_enabled"`
	PaymentEnabled     bool        `json:"payment_enabled"`
	Form               []Form      `json:"form"`
	AuthenticationType string      `json:"authentication_type"`
}

type Form struct {
	Label     string      `json:"label"`
	Type      string      `json:"type"`
	IsNum     string      `json:"isNum"`
	MaxLength interface{} `json:"maxLength"`
}
