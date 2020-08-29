package main

import (
	"net/http"
	"os"
	"time"

	"github.com/Seklfreak/bridgeapi-lunchmoney-sync/bridgeapi"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	bridgeClient, err := bridgeapi.NewClient(&http.Client{}, &bridgeapi.Auth{
		ClientID:     os.Getenv("BRIDGEAPI_CLIENT_ID"),
		ClientSecret: os.Getenv("BRIDGEAPI_CLIENT_SECRET"),
		Email:        os.Getenv("BRIDGEAPI_EMAIL"),
		Password:     os.Getenv("BRIDGEAPI_PASSWORD"),
	})
	if err != nil {
		logger.Fatal("failure initialising bridge client", zap.Error(err))
	}

	// fetch updated in last seven days
	transactions, err := bridgeClient.FetchTransactionsUpdated(time.Now().Add(-7 * 24 * time.Hour))
	if err != nil {
		logger.Fatal("failure fetching transactions from bridge", zap.Error(err))
	}

	logger.Info("received transactions", zap.Int("amount", len(transactions)))
}
