package smodels

import "github.com/shopspring/decimal"

type (
	Pie struct {
		Parts []PiePart       `json:"parts"`
		Total decimal.Decimal `json:"total"`
	}
	PiePart struct {
		Label string          `json:"label"`
		Title string          `json:"title"`
		Value decimal.Decimal `json:"value"`
	}
)
