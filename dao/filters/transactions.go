package filters

type Transactions struct {
	Height uint64 `schema:"height"`
	Limit  uint64 `schema:"limit"`
	Offset uint64 `schema:"offset"`
}
