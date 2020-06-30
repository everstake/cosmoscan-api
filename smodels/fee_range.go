package smodels

import "github.com/shopspring/decimal"

type FeeRange struct {
	From       decimal.Decimal     `json:"from"`
	To         decimal.Decimal     `json:"to"`
	Validators []FeeRangeValidator `json:"validators"`
}

type FeeRangeValidator struct {
	Validator string          `json:"validator"`
	Fee       decimal.Decimal `json:"fee"`
}
