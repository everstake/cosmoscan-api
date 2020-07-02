package dmodels

import "github.com/shopspring/decimal"

type ValidatorDelegator struct {
	Delegator string          `schema:"delegator"`
	Amount    decimal.Decimal `schema:"amount"`
	Since     Time            `schema:"since"`
	Delta     decimal.Decimal `schema:"delta"`
}
