package filters

type Blocks struct {
	Limit  uint64
	Offset uint64
}

type BlocksProposed struct {
	Proposers []string
}
