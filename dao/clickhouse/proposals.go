package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (db DB) CreateProposal(proposals []dmodels.Proposal) error {
	if len(proposals) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.ProposalsTable).Columns("pro_id", "pro_init_deposit", "pro_proposer", "pro_content", "pro_created_at")
	for _, proposal := range proposals {
		if proposal.ID == "" {
			return fmt.Errorf("field ID can not be empty")
		}
		if proposal.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(proposal.ID, proposal.InitDeposit, proposal.Proposer, proposal.Content, proposal.CreatedAt)
	}
	return db.Insert(q)
}
