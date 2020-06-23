package filters

type VotingPower struct {
	TimeRange
	Delegators []string
	Validators []string
}
