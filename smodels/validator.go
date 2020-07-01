package smodels

import "github.com/shopspring/decimal"

type Validator struct {
	Title           string          `json:"title"`
	Website         string          `json:"website"`
	OperatorAddress string          `json:"operator_address"`
	AccAddress      string          `json:"acc_address"`
	ConsAddress     string          `json:"cons_address"`
	PercentPower    decimal.Decimal `json:"percent_power"`
	Power           decimal.Decimal `json:"power"`
	SelfStake       decimal.Decimal `json:"self_stake"`
	Fee             decimal.Decimal `json:"fee"`
	BlocksProposed  uint64          `json:"blocks_proposed"`
	Delegators      uint64          `json:"delegators"`
	Power24Change   decimal.Decimal `json:"power_24_change"`
	GovernanceVotes uint64          `json:"governance_votes"`
}
