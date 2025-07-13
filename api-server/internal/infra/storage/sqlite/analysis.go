package sqlite

import (
	"context"
	"database/sql"
	"log"

	"api-server/domain"

	_ "github.com/mattn/go-sqlite3" // Import SQLite driver
)

type analysisStorage struct {
	DB  *sql.DB
	log *log.Logger
}

func NewAnalysisStorage(db *sql.DB, log *log.Logger) (*analysisStorage, error) {
	return &analysisStorage{
		DB:  db,
		log: log,
	}, nil
}

func (s *analysisStorage) CreateAnalysis(ctx context.Context, analysis *domain.AwesomeAPIResponse) error {
	query := `INSERT INTO dolar_quotes (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := s.DB.ExecContext(ctx, query, analysis.DolarQuote.Code, analysis.DolarQuote.CodeIn, analysis.DolarQuote.Name,
		analysis.DolarQuote.High, analysis.DolarQuote.Low, analysis.DolarQuote.VarBid,
		analysis.DolarQuote.PercentChange, analysis.DolarQuote.Bid, analysis.DolarQuote.Ask,
		analysis.DolarQuote.Timestamp, analysis.DolarQuote.CreateDate)
	if err != nil {
		s.log.Printf("erro ao salvar a cotação no banco de dados: %v", err)
		return err
	}

	return nil
}
