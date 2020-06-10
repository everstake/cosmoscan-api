package filters

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
	"time"
)

const (
	AggByHour  = "hour"
	AggByDay   = "day"
	AggByWeek  = "week"
	AggByMonth = "month"
)

type Agg struct {
	By   string       `schema:"by"`
	From dmodels.Time `schema:"from"`
	To   dmodels.Time `schema:"to"`
}

var aggLimits = map[string]struct {
	defaultRange time.Duration
	maxRange     time.Duration
}{
	AggByHour: {
		defaultRange: time.Hour * 24,
		maxRange:     time.Hour * 24 * 7,
	},
	AggByDay: {
		defaultRange: time.Hour * 24 * 30,
		maxRange:     time.Hour * 24 * 30 * 2,
	},
	AggByWeek: {
		defaultRange: time.Hour * 24 * 40,
		maxRange:     time.Hour * 24 * 40 * 3,
	},
	AggByMonth: {
		defaultRange: time.Hour * 24 * 365,
		maxRange:     time.Hour * 24 * 365 * 2,
	},
}

func (agg *Agg) Validate() error {
	limit, ok := aggLimits[agg.By]
	if !ok {
		return fmt.Errorf("not found `by` param")
	}
	if agg.From.IsZero() {
		agg.From = dmodels.NewTime(time.Now().Add(-limit.defaultRange))
		agg.To = dmodels.NewTime(time.Now())
	} else {
		d := agg.To.Sub(agg.From.Time)
		if d > limit.maxRange {
			return fmt.Errorf("over max limit range")
		}
	}
	return nil
}

func (agg *Agg) AggFunc() string {
	switch agg.By {
	case AggByHour:
		return "toStartOfHour"
	case AggByDay:
		return "toStartOfDay"
	case AggByWeek:
		return "toStartOfWeek"
	case AggByMonth:
		return "toStartOfMonth"
	default:
		return "toStartOfDay"
	}
}

func (agg *Agg) BuildQuery(aggValue string, timeColumn string, table string) squirrel.SelectBuilder {
	q := squirrel.Select(
		fmt.Sprintf("%s as value", aggValue),
		fmt.Sprintf("toDateTime(%s(%s)) AS time", agg.AggFunc(), timeColumn),
	).From(table).
		GroupBy("time").
		OrderBy("time")
	if !agg.From.IsZero() {
		q = q.Where(squirrel.GtOrEq{timeColumn: agg.From.Time})
	}
	if !agg.To.IsZero() {
		q = q.Where(squirrel.LtOrEq{timeColumn: agg.To.Time})
	}
	return q
}
