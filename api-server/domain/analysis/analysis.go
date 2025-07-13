package analysis

import (
	"api-server/domain"
	"context"
	"errors"
	"log"
	"time"
)

type analysisService struct {
	analysisStorage  domain.AnalysisStorage
	awesomeAPIClient domain.AwesomeAPIClient
	log              *log.Logger
}

func NewAnalysisService(analysisStorage domain.AnalysisStorage, awesomeAPIClient domain.AwesomeAPIClient, log *log.Logger) *analysisService {

	return &analysisService{
		analysisStorage:  analysisStorage,
		awesomeAPIClient: awesomeAPIClient,
		log:              log,
	}
}

func (s *analysisService) RunAnalysis(c context.Context) (string, error) {
	// Timeout de 200ms para chamada da API de cotação
	apiCtx, apiCancel := context.WithTimeout(c, 200*time.Millisecond)
	defer apiCancel()

	dolarQuote, err := s.awesomeAPIClient.GetDolarQuote(apiCtx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.log.Printf("timeout para obter cotação do dolar da AwesomeAPI: %v", err)
		} else {
			s.log.Printf("erro ao obter a cotação do dolar: %v", err)
		}
		return "", err
	}

	// Timeout de 10ms para persistência no banco
	dbCtx, dbCancel := context.WithTimeout(c, 10*time.Millisecond)
	defer dbCancel()

	err = s.analysisStorage.CreateAnalysis(dbCtx, dolarQuote)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			s.log.Printf("timeout para persistir cotação do dolar no banco de dados: %v", err)
		} else {
			s.log.Printf("erro ao salvar a cotação do dolar: %v", err)
		}
		return "", err
	}

	return dolarQuote.DolarQuote.VarBid, nil
}
