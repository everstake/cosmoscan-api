package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const ValidatorsTable = "validators"

type Validator struct {
	Address         string          `db:"val_address"`
	OperatorAddress string          `db:"val_operator_address"`
	ConsAddress     string          `db:"val_cons_address"`
	ConsPubKey      string          `db:"val_cons_pub_key"`
	Name            string          `db:"val_name"`
	Description     string          `db:"val_description"`
	Commission      decimal.Decimal `db:"val_commission"`
	MinCommission   decimal.Decimal `db:"val_min_commission"`
	MaxCommission   decimal.Decimal `db:"val_max_commission"`
	SelfDelegations decimal.Decimal `db:"val_self_delegations"`
	Delegations     decimal.Decimal `db:"val_delegations"`
	VotingPower     decimal.Decimal `db:"val_voting_power"`
	Website         string          `db:"val_website"`
	Jailed          bool            `db:"val_jailed"`
	CreatedAt       time.Time       `db:"val_created_at"`
}
