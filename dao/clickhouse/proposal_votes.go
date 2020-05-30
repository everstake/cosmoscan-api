package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (db DB) CreateProposalVotes(votes []dmodels.ProposalVote) error {
	if len(votes) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.ProposalVotesTable).Columns("prv_id", "prv_proposal_id", "prv_voter", "prv_option", "prv_created_at")
	for _, vote := range votes {
		if vote.ID == "" {
			return fmt.Errorf("field ID can not be empty")
		}
		if vote.ProposalID == 0 {
			return fmt.Errorf("field ProposalID can not be zero")
		}
		if vote.Voter == "" {
			return fmt.Errorf("field Voter can not be zero")
		}
		if vote.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(vote.ID, vote.ProposalID, vote.Voter, vote.Option, vote.CreatedAt)
	}
	return db.Insert(q)
}
