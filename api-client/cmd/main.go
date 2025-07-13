package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CotacaoResponse struct {
	DolarBid string `json:"dolar_bid"`
	Message  string `json:"message"`
}

func main() {
	logger := log.New(os.Stdout, "api-client", log.LstdFlags)
	// Timeout de 300ms para chamada de retorno do servidor
	apiCtx, apiCancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer apiCancel()

	URL := "http://localhost:8080/cotacao"

	req, err := http.NewRequestWithContext(apiCtx, http.MethodGet, URL, nil)
	if err != nil {
		logger.Printf("erro ao criar a requisição para o servidor: %s", err.Error())
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Printf("erro ao chamar o servidor: %s", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Printf("erro ao ler resposta do servidor: %s", err.Error())
		return
	}

	var cotacao CotacaoResponse
	if err := json.Unmarshal(body, &cotacao); err != nil {
		logger.Printf("erro ao decodificar resposta JSON: %s", err.Error())
		return
	}

	file, err := os.Create("cotacao.txt")
	if err != nil {
		logger.Printf("erro ao criar arquivo cotacao.txt: %s", err.Error())
		return
	}
	defer file.Close()

	_, err = file.WriteString("Dólar: " + cotacao.DolarBid)
	if err != nil {
		logger.Printf("erro ao escrever no arquivo cotacao.txt: %s", err.Error())
		return
	}

	logger.Println("Cotação salva em cotacao.txt")
}
