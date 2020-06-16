package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const BalanceUpdatesTable = "balance_updates"

type BalanceUpdate struct {
	ID        string          `db:"bau_id"`
	Address   string          `db:"bau_address"`
	Stake     decimal.Decimal `db:"bau_stake"`
	Balance   decimal.Decimal `db:"bau_balance"`
	Unbonding decimal.Decimal `db:"bau_unbonding"`
	CreatedAt time.Time       `db:"bau_created_at"`
}
