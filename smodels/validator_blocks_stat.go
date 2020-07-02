package smodels

import "github.com/shopspring/decimal"

type ValidatorBlocksStat struct {
	Proposed          uint64          `json:"proposed"`
	MissedValidations uint64          `json:"missed_validations"`
	Revenue           decimal.Decimal `json:"revenue"`
}
