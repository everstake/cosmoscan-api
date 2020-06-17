package filters

type Proposals struct {
	ID     []uint64 `schema:"id"`
	Limit  uint64   `schema:"limit"`
	Offset uint64   `schema:"offset"`
}
