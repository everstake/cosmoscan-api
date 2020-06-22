package dmodels

const ProposalVotesTable = "proposal_votes"

type ProposalVote struct {
	ID         string `db:"prv_id" json:"-"`
	ProposalID uint64 `db:"prv_proposal_id" json:"proposal_id"`
	Voter      string `db:"prv_voter" json:"voter"`
	TxHash     string `db:"prv_tx_hash" json:"tx_hash"`
	Option     string `db:"prv_option" json:"option"`
	CreatedAt  Time   `db:"prv_created_at" json:"created_at"`
}
