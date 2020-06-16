package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const AccountsTable = "accounts"

type Account struct {
	Address   string          `db:"acc_address"`
	Balance   decimal.Decimal `db:"acc_balance"`
	Stake     decimal.Decimal `db:"acc_stake"`
	Unbonding decimal.Decimal `db:"acc_unbonding"`
	CreatedAt time.Time       `db:"acc_created_at"`
}
