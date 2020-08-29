package main

import (
	"strings"

	"github.com/Seklfreak/bridgeapi-lunchmoney-sync/bridgeapi"
	"github.com/Seklfreak/bridgeapi-lunchmoney-sync/lunchmoney"
)

type accountWrapped struct {
	Account *bridgeapi.Account
	Bank    *bridgeapi.Bank
}

func matchToAsset(assets []*lunchmoney.Asset, accounts []*accountWrapped, transaction *bridgeapi.Transaction) int {
	account := findAccountForTrx(accounts, transaction)
	if account == nil {
		return 0
	}

	for _, asset := range assets {
		if strings.Contains(account.Account.Name, asset.Name) &&
			strings.Contains(account.Bank.Name, asset.Name) {
			return asset.ID
		}
	}

	return 0
}

func findAccountForTrx(accounts []*accountWrapped, transaction *bridgeapi.Transaction) *accountWrapped {
	for _, account := range accounts {
		if account.Account.ID == transaction.Account.ID {
			return account
		}
	}

	return nil
}
