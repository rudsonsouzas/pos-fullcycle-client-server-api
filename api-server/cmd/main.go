package main

import (
	"os"
	"time"

	"os/signal"
	"syscall"

	"api-server/domain/analysis"
	"api-server/internal/infra/client"
	"api-server/internal/infra/server/http"
	"api-server/internal/infra/storage/sqlite"
	"api-server/pkg/env"
	httpclient "api-server/pkg/http_client"
	"log"
)

const (
	envApplicationPort = "APP_PORT"
	envAwesomeAPIURL   = "AWESOMEAPI_URL"

	defaultApplicationPort = "8080"
	defaultAwesomeAPIURL   = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
)

var (
	version, build, date string
)

func main() {
	logger := log.New(os.Stdout, "api-server - ", log.LstdFlags)

	logger.Printf("API Cotacao Dolar - version:%s; build:%s; date:%s", version, build, date)

	awesomeAPIClient := client.NewAwesomeAPIClient(httpclient.NewHTTPClient(60*time.Second), getAwesomeAPIURL(), logger)
	/*
	* Storages
	 */
	sqlDB, err := sqlite.ConnectDB()
	if err != nil {
		logger.Fatalf("error connecting database sqlite: %q", err)
		return
	}
	analisysClient, err := sqlite.NewAnalysisStorage(sqlDB, logger)
	if err != nil {
		logger.Fatalf("error creating database client: %q", err)
		return
	}

	analysisService := analysis.NewAnalysisService(analisysClient, awesomeAPIClient, logger)

	handler := http.NewHandler(analysisService, logger)

	/*
	 * Server...
	 */
	server := http.New(getApplicationPort(), handler, logger)
	server.ListenAndServe()

	/*
	 * Graceful shutdown...
	 */
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
	server.Shutdown()
}

func getApplicationPort() string {
	return env.GetString(envApplicationPort, defaultApplicationPort)
}

func getAwesomeAPIURL() string {
	return env.GetString(envAwesomeAPIURL, defaultAwesomeAPIURL)
}
