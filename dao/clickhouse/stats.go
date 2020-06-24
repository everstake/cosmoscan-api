package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/smodels"
)

func (db DB) CreateStats(stats []dmodels.Stat) (err error) {
	if len(stats) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.StatsTable).Columns("stt_id", "stt_title", "stt_value", "stt_created_at")
	for _, stat := range stats {
		if stat.ID == "" {
			return fmt.Errorf("field ProposalID can not be empty")
		}
		if stat.Title == "" {
			return fmt.Errorf("field Title can not be empty")
		}
		if stat.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(stat.ID, stat.Title, stat.Value, stat.CreatedAt)
	}
	return db.Insert(q)
}

func (db DB) GetStats(filter filters.Stats) (stats []dmodels.Stat, err error) {
	q := squirrel.Select("*").From(dmodels.StatsTable).OrderBy("stt_created_at")
	if !filter.From.IsZero() {
		q = q.Where(squirrel.GtOrEq{"stt_created_at": filter.From})
	}
	if !filter.To.IsZero() {
		q = q.Where(squirrel.LtOrEq{"stt_created_at": filter.To})
	}
	if len(filter.Titles) != 0 {
		q = q.Where(squirrel.Eq{"stt_title": filter.Titles})
	}
	err = db.Find(&stats, q)
	return stats, err
}

func (db DB) GetAggValidators33Power(filter filters.Agg) (items []smodels.AggItem, err error) {
	q := filter.BuildQuery("max(stt_value)", "stt_created_at", dmodels.StatsTable).
		Where(squirrel.Eq{"stt_title": dmodels.StatsValidatorsWith33Power})
	err = db.Find(&items, q)
	return items, err
}

func (db DB) GetAggWhaleAccounts(filter filters.Agg) (items []smodels.AggItem, err error) {
	q := filter.BuildQuery("max(stt_value)", "stt_created_at", dmodels.StatsTable).
		Where(squirrel.Eq{"stt_title": dmodels.StatsTotalWhaleAccounts})
	err = db.Find(&items, q)
	return items, err
}
