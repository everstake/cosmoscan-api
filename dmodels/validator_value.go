package dmodels

type ValidatorValue struct {
	Validator string `db:"validator"`
	Value     uint64 `db:"value"`
}
