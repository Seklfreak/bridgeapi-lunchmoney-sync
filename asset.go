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

func machTransactionToAsset(assets []*lunchmoney.Asset, accounts []*accountWrapped, transaction *bridgeapi.Transaction) *lunchmoney.Asset {
	account := findAccountForTrx(accounts, transaction)
	if account == nil {
		return nil
	}

	return matchAccountToAsset(assets, account)
}

func matchAccountToAsset(assets []*lunchmoney.Asset, account *accountWrapped) *lunchmoney.Asset {
	// try strict mapping
	for _, asset := range assets {
		if (strings.Contains(account.Account.Name, asset.Name)) &&
			(strings.Contains(account.Bank.Name, asset.InstitutionName)) {
			return asset
		}
	}

	// try more relaxed mapping
	for _, asset := range assets {
		if (strings.Contains(account.Account.Name, asset.Name) || strings.Contains(account.Account.Name, asset.InstitutionName)) &&
			(strings.Contains(account.Bank.Name, asset.Name) || strings.Contains(account.Bank.Name, asset.InstitutionName)) {
			return asset
		}
	}

	return nil
}

func findAccountForTrx(accounts []*accountWrapped, transaction *bridgeapi.Transaction) *accountWrapped {
	for _, account := range accounts {
		if account.Account.ID == transaction.Account.ID {
			return account
		}
	}

	return nil
}
