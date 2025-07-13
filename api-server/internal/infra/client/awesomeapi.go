package client

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"api-server/domain"
	httpclient "api-server/pkg/http_client"

	"github.com/cenkalti/backoff"
)

type AwesomeAPIClient struct {
	httpClient httpclient.HTTPClient
	baseURL    string
	log        *log.Logger
}

// type AwesomeAPIResponse struct {
// 	Result map[string]any
// 	Error  error
// }

type AwesomeAPI interface {
	GetDolarQuote(c context.Context) (*domain.AwesomeAPIResponse, error)
}

func NewAwesomeAPIClient(httpClient httpclient.HTTPClient, baseURL string, log *log.Logger) *AwesomeAPIClient {
	return &AwesomeAPIClient{
		httpClient: httpClient,
		baseURL:    baseURL,
		log:        log,
	}
}

func (awc *AwesomeAPIClient) GetDolarQuote(ctx context.Context) (*domain.AwesomeAPIResponse, error) {
	response := &domain.AwesomeAPIResponse{}

	ebo := backoff.NewExponentialBackOff()
	ebo.MaxInterval = 10 * time.Second

	if err := backoff.Retry(func() error {

		URL := awc.baseURL

		req, err := http.NewRequest(http.MethodGet, URL, nil)
		if err != nil {
			awc.log.Printf("erro ao criar a requisição para a AwesomeAPI: %s", err.Error())
			return err
		}

		req.Header.Add("Content-Type", "application/json")

		res, err := awc.httpClient.Do(req)
		if err != nil {
			awc.log.Printf("erro ao realizar a requisição para a AwesomeAPI: %s", err.Error())
			return err
		}
		defer func() {
			err = res.Body.Close()
			if err != nil {
				awc.log.Printf("erro ao encerrar o response body: %s", err.Error())
				return
			}
		}()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			awc.log.Printf("erro ao ler o response body: %s", err.Error())
			return err
		}

		if res.StatusCode != 200 {
			awc.log.Printf("Awesome API status code %d: %s", res.StatusCode, bodyBytes)
			return err
		}

		if err := json.Unmarshal(bodyBytes, &response); err != nil {
			awc.log.Printf("erro ao traduzir o response Body para o padrão esperado: %s", err.Error())
			return err
		}

		return nil

	}, backoff.WithContext(backoff.WithMaxRetries(ebo, uint64(5)), context.Background())); err != nil {
		awc.log.Printf("erro ao requisitar a cotação do dolar da Awesome API: %s", err.Error())
		return response, err
	}

	return response, nil
}
