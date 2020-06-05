package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const BalancesUpdatesTable = "balance_updates"

type BalanceUpdate struct {
	Address   string          `db:"bau_address"`
	Balance   decimal.Decimal `db:"bau_balance"`
	Delta     decimal.Decimal `db:"bau_delta"`
	CreatedAt time.Time       `db:"bau_created_at"`
}
