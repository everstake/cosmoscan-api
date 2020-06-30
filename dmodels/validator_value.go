package dmodels

type ValidatorValue struct {
	Validator string `db:"validator" json:"validator"`
	Value     uint64 `db:"value" json:"value"`
}
