package dmodels

import (
	"github.com/shopspring/decimal"
	"time"
)

const ProposalDepositsTable = "proposal_deposits"

type ProposalDeposit struct {
	ID         string          `db:"prd_id"`
	ProposalID uint64          `db:"prd_proposal_id"`
	Depositor  string          `db:"prd_depositor"`
	Amount     decimal.Decimal `db:"prd_amount"`
	CreatedAt  time.Time       `db:"prd_created_at"`
}
