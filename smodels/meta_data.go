package smodels

import "github.com/shopspring/decimal"

type MetaData struct {
	Height          uint64          `json:"height"`
	LatestValidator string          `json:"latest_validator"`
	LatestProposal  string          `json:"latest_proposal"`
	ValidatorAvgFee decimal.Decimal `json:"validator_avg_fee"`
	BlockTime       float64         `json:"block_time"`
	CurrentPrice    decimal.Decimal `json:"current_price"`
}
