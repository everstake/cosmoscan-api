package filters

import (
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
)

type TimeRange struct {
	From dmodels.Time `schema:"from"`
	To   dmodels.Time `schema:"to"`
}

func (filter *TimeRange) Query(timeColumn string, q squirrel.SelectBuilder) squirrel.SelectBuilder {
	if !filter.From.IsZero() {
		q = q.Where(squirrel.GtOrEq{timeColumn: filter.From.Time})
	}
	if !filter.To.IsZero() {
		q = q.Where(squirrel.LtOrEq{timeColumn: filter.To.Time})
	}
	return q
}
