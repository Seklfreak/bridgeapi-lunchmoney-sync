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

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	bridgeClient, err := bridgeapi.NewClient(&http.Client{}, &bridgeapi.Auth{
		ClientID:     os.Getenv("BRIDGEAPI_CLIENT_ID"),
		ClientSecret: os.Getenv("BRIDGEAPI_CLIENT_SECRET"),
		Email:        os.Getenv("BRIDGEAPI_EMAIL"),
		Password:     os.Getenv("BRIDGEAPI_PASSWORD"),
	})
	if err != nil {
		logger.Fatal("failure initialising bridge client", zap.Error(err))
	}

	lunchmoneyClient := lunchmoney.NewClient(&http.Client{}, os.Getenv("LUNCHMONEY_ACCESS_TOKEN"))

	// fetch updated in last seven days
	transactions, err := bridgeClient.FetchTransactionsUpdated(time.Now().Add(-7 * 24 * time.Hour))
	if err != nil {
		logger.Fatal("failure fetching transactions from bridge", zap.Error(err))
	}

	logger.Info("received transactions", zap.Int("amount", len(transactions)))

	var convertedTrxs []*lunchmoney.Transaction
	var notes string
	for _, trx := range transactions {
		notes = ""
		if !strings.EqualFold(trx.Description, trx.RawDescription) {
			notes = trx.RawDescription
		}

		convertedTrxs = append(convertedTrxs, &lunchmoney.Transaction{
			Date:        trx.Date,
			Amount:      trx.Amount,
			CategoryID:  0,
			Payee:       trx.Description,
			Currency:    strings.ToLower(trx.CurrencyCode),
			AssetID:     0, // TODO
			RecurringID: 0,
			Notes:       notes,
			Status:      "",
			ExternalID:  strconv.FormatInt(trx.ID, 10),
			Tags:        []string{"bridgeapi-lunchmoney-sync"},
		})

		// TODO
		if len(convertedTrxs) >= 10 {
			break
		}
	}

	inserted, err := lunchmoneyClient.InsertTransactions(ctx, convertedTrxs)
	if err != nil {
		logger.Fatal("failure inserting transactions to lunchmoney", zap.Error(err))
	}

	logger.Info("inserted transactions", zap.Int("amount", inserted))
}
