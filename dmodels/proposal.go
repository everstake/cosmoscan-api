package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const ProposalsTable = "proposals"

type Proposal struct {
	ID          uint64          `db:"pro_id"`
	Title       string          `db:"pro_title"`
	Description string          `db:"pro_description"`
	Recipient   string          `db:"pro_recipient"`
	Amount      decimal.Decimal `db:"pro_amount"`
	InitDeposit decimal.Decimal `db:"pro_init_deposit"`
	Proposer    string          `db:"pro_proposer"`
	CreatedAt   time.Time       `db:"pro_created_at"`
}
