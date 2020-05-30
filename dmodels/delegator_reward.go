package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const DelegatorRewardsTable = "rewards"

type DelegatorReward struct {
	ID        string          `db:"der_id"`
	TxHash    string          `db:"der_tx_hash"`
	Delegator string          `db:"der_delegator"`
	Validator string          `db:"der_validator"`
	Amount    decimal.Decimal `db:"der_amount"`
	CreatedAt time.Time       `db:"der_created_at"`
}
