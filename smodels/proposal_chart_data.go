package smodels

import "github.com/shopspring/decimal"

type ProposalChartData struct {
	ProposalID        uint64          `json:"proposal_id"`
	VotersTotal       uint64          `json:"voters_total"`
	ValidatorsTotal   uint64          `json:"validators_total"`
	Turnout           decimal.Decimal `json:"turnout"`
	YesPercent        decimal.Decimal `json:"yes_percent"`
	NoPercent         decimal.Decimal `json:"no_percent"`
	NoWithVetoPercent decimal.Decimal `json:"no_with_veto_percent"`
	AbstainPercent    decimal.Decimal `json:"abstain_percent"`
}
