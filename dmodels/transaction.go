package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const TransactionsTable = "transactions"

type Transaction struct {
	Hash      string          `db:"trn_hash"`
	Status    bool            `db:"trn_status"`
	Height    uint64          `db:"trn_height"`
	Messages  uint64          `db:"trn_messages"`
	Fee       decimal.Decimal `db:"trn_fee"`
	GasUsed   uint64          `db:"trn_gas_used"`
	GasWanted uint64          `db:"trn_gas_wanted"`
	CreatedAt time.Time       `db:"trn_created_at"`
}
