package dmodels

import (
	"github.com/shopspring/decimal"
)

const ProposalDepositsTable = "proposal_deposits"

type ProposalDeposit struct {
	ID         string          `db:"prd_id" json:"-"`
	ProposalID uint64          `db:"prd_proposal_id" json:"proposal_id"`
	Depositor  string          `db:"prd_depositor" json:"depositor"`
	Amount     decimal.Decimal `db:"prd_amount" json:"amount"`
	CreatedAt  Time            `db:"prd_created_at" json:"created_at"`
}
