package dmodels

import (
	"time"
)

const ProposalVotesTable = "proposal_votes"

type ProposalVote struct {
	ID         string    `db:"prv_id"`
	ProposalID uint64    `db:"prv_proposal_id"`
	Voter      string    `db:"prv_voter"`
	Option     string    `db:"option"`
	CreatedAt  time.Time `db:"prv_created_at"`
}
