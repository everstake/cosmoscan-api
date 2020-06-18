package dmodels

import (
	"time"
)

const ProposalVotesTable = "proposal_votes"

type ProposalVote struct {
	ID         string    `db:"prv_id"`
	ProposalID uint64    `db:"prv_proposal_id"`
	Voter      string    `db:"prv_voter"`
	TxHash     string    `db:"prv_tx_hash"`
	Option     string    `db:"prv_option"`
	CreatedAt  time.Time `db:"prv_created_at"`
}
