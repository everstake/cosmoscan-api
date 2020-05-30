package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const ProposalsTable = "proposals"

type Proposal struct {
	ID          string          `db:"pro_id"`
	InitDeposit decimal.Decimal `db:"pro_init_deposit"`
	Proposer    string          `db:"pro_proposer"`
	Content     string          `db:"pro_content"`
	CreatedAt   time.Time       `db:"pro_created_at"`
}
