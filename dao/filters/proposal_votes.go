package filters

type ProposalVotes struct {
	ProposalID []uint64 `schema:"proposal_id"`
	Voters     []string `schema:"voters"`
	Limit      uint64   `schema:"limit"`
	Offset     uint64   `schema:"offset"`
}
