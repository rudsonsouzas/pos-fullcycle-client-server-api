package domain

import "context"

type AwesomeAPIResponse struct {
	DolarQuote DolarQuoteInfo `json:"USDBRL"`
}

type DolarQuoteInfo struct {
	Code          string `json:"code"`
	CodeIn        string `json:"codein"`
	Name          string `json:"name"`
	High          string `json:"high"`
	Low           string `json:"low"`
	VarBid        string `json:"varBid"`
	PercentChange string `json:"pctChange"`
	Bid           string `json:"bid"`
	Ask           string `json:"ask"`
	Timestamp     string `json:"timestamp"`
	CreateDate    string `json:"create_date"`
}

type AwesomeAPIClient interface {
	GetDolarQuote(ctx context.Context) (*AwesomeAPIResponse, error)
}
