package dmodels

import "github.com/shopspring/decimal"

type ValidatorDelegator struct {
	Delegator string          `json:"delegator"`
	Amount    decimal.Decimal `json:"amount"`
	Since     Time            `json:"since"`
	Delta     decimal.Decimal `json:"delta"`
}
