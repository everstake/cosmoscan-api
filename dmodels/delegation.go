package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const DelegationsTable = "delegations"

type Delegation struct {
	ID        string          `db:"dlg_id"`
	TxHash    string          `db:"dlg_tx_hash"`
	Delegator string          `db:"dlg_delegator"`
	Validator string          `db:"dlg_validator"`
	Amount    decimal.Decimal `db:"dlg_amount"`
	CreatedAt time.Time       `db:"dlg_created_at"`
}
