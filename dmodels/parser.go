package dmodels

const ParsersTable = "parsers"

type Parser struct {
	ID     uint64 `db:"par_id"`
	Title  string `db:"par_title"`
	Height uint64 `db:"par_height"`
}
