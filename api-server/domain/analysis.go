package domain

import "context"

type AnalysisStorage interface {
	CreateAnalysis(ctx context.Context, analysis *AwesomeAPIResponse) error
}

type AnalysisService interface {
	RunAnalysis(c context.Context) (string, error)
}
