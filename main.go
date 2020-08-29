package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Seklfreak/bridgeapi-lunchmoney-sync/bridgeapi"
	"github.com/Seklfreak/bridgeapi-lunchmoney-sync/lunchmoney"
	"go.uber.org/zap"
)

var tags = []string{"bridgeapi-lunchmoney-sync"}

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	bridgeClient, err := bridgeapi.NewClient(ctx, &http.Client{}, &bridgeapi.Auth{
		ClientID:     os.Getenv("BRIDGEAPI_CLIENT_ID"),
		ClientSecret: os.Getenv("BRIDGEAPI_CLIENT_SECRET"),
		Email:        os.Getenv("BRIDGEAPI_EMAIL"),
		Password:     os.Getenv("BRIDGEAPI_PASSWORD"),
	})
	if err != nil {
		logger.Fatal("failure initialising bridge client", zap.Error(err))
	}

	lunchmoneyClient := lunchmoney.NewClient(&http.Client{}, os.Getenv("LUNCHMONEY_ACCESS_TOKEN"))

	// fetch lunchmoney assets
	assets, err := lunchmoneyClient.GetAssets(ctx)
	if err != nil {
		logger.Fatal("failure fetching lunchmoney assets", zap.Error(err))
	}

	// fetch bridge accounts
	accounts, err := bridgeClient.FetchAccounts(ctx)
	if err != nil {
		logger.Fatal("failure fetching bridge accounts", zap.Error(err))
	}

	var accountsWrapped []*accountWrapped

	// fetch bridge banks
	for _, account := range accounts {
		bank, err := bridgeClient.FetchBank(ctx, account.Bank.ID)
		if err != nil {
			logger.Fatal("failure fetching bridge bank", zap.Error(err), zap.Int("bank_id", account.Bank.ID))
		}

		accountsWrapped = append(accountsWrapped, &accountWrapped{Account: account, Bank: bank})
	}

	// fetch updated in last seven days
	transactions, err := bridgeClient.FetchTransactionsUpdated(ctx, time.Now().Add(-7*24*time.Hour))
	if err != nil {
		logger.Fatal("failure fetching transactions from bridge", zap.Error(err))
	}

	logger.Info("received transactions", zap.Int("amount", len(transactions)))

	// convert to lunchmoney transactions
	var convertedTrxs []*lunchmoney.Transaction
	var notes string
	var assetID int
	for _, trx := range transactions {
		assetID = matchToAsset(assets, accountsWrapped, trx)
		if assetID == 0 {
			logger.Warn("failure matching transaction to asset", zap.Any("trx", trx))
			continue
		}

		notes = ""
		if !strings.EqualFold(trx.Description, trx.RawDescription) {
			notes = trx.RawDescription
		}

		convertedTrxs = append(convertedTrxs, &lunchmoney.Transaction{
			Date:       trx.Date,
			Amount:     trx.Amount,
			Payee:      trx.Description,
			Currency:   strings.ToLower(trx.CurrencyCode),
			AssetID:    assetID,
			Notes:      notes,
			ExternalID: strconv.FormatInt(trx.ID, 10),
			Tags:       tags,
		})
	}

	// send to lunchmoney
	inserted, err := lunchmoneyClient.InsertTransactions(ctx, convertedTrxs)
	if err != nil {
		logger.Fatal("failure inserting transactions to lunchmoney", zap.Error(err))
	}

	logger.Info("inserted transactions", zap.Int("amount", inserted))
}
