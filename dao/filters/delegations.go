package filters

type Delegators struct {
	TimeRange
	Validators []string `schema:"validators"`
}

type DelegationsAgg struct {
	Agg
	Validators []string `schema:"validators"`
}

type ValidatorDelegators struct {
	Validator string `json:"-"`
	Limit     uint64 `schema:"limit"`
	Offset    uint64 `schema:"offset"`
}
