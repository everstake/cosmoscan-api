package filters

type Delegators struct {
	TimeRange
	Validators []string `schema:"validators"`
}
