package smodels

import "github.com/shopspring/decimal"

type Balance struct {
	SelfDelegated  decimal.Decimal `json:"self_delegated"`
	OtherDelegated decimal.Decimal `json:"other_delegated"`
	Available      decimal.Decimal `json:"available"`
}
