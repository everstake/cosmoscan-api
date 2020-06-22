package filters

type ProposalVotes struct {
	ProposalID []uint64 `schema:"proposal_id"`
	Limit      uint64   `schema:"limit"`
	Offset     uint64   `schema:"offset"`
}
