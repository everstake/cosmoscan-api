package filters

type Blocks struct {
	Limit  uint64 `schema:"limit"`
	Offset uint64 `schema:"offset"`
}

type BlocksProposed struct {
	Proposers []string
}
