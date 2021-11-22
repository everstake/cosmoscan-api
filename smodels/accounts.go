package smodels

import "github.com/shopspring/decimal"

type Account struct {
	Address     string          `json:"address"`
	Balance     decimal.Decimal `json:"balance"`
	Delegated   decimal.Decimal `json:"delegated"`
	Unbonding   decimal.Decimal `json:"unbonding"`
	StakeReward decimal.Decimal `json:"stake_reward"`
}
