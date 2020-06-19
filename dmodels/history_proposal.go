package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const HistoryProposalsTable = "history_proposals"

type HistoryProposal struct {
	ID          uint64          `db:"hpr_id"`
	TxHash      string          `db:"hpr_tx_hash"`
	Title       string          `db:"hpr_title"`
	Description string          `db:"hpr_description"`
	Recipient   string          `db:"hpr_recipient"`
	Amount      decimal.Decimal `db:"hpr_amount"`
	InitDeposit decimal.Decimal `db:"hpr_init_deposit"`
	Proposer    string          `db:"hpr_proposer"`
	CreatedAt   time.Time       `db:"hpr_created_at"`
}
