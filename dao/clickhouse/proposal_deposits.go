package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (db DB) CreateProposalDeposits(deposits []dmodels.ProposalDeposit) error {
	if len(deposits) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.ProposalDepositsTable).Columns("prd_id", "prd_proposal_id", "prd_depositor", "prd_amount", "prd_created_at")
	for _, deposit := range deposits {
		if deposit.ID == "" {
			return fmt.Errorf("field ID can not be empty")
		}
		if deposit.ProposalID == 0 {
			return fmt.Errorf("field ProposalID can not be zero")
		}
		if deposit.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(deposit.ID, deposit.ProposalID, deposit.Depositor, deposit.Amount, deposit.CreatedAt)
	}
	return db.Insert(q)
}
