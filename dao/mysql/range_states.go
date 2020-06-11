package mysql

import (
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (m DB) GetRangeStates(titles []string) (items []dmodels.RangeState, err error) {
	q := squirrel.Select("*").From(dmodels.RangeStatesTable)
	if len(titles) > 0 {
		q = q.Where(squirrel.Eq{"rst_title": titles})
	}
	err = m.find(&items, q)
	return items, err
}

func (m DB) UpdateRangeState(item dmodels.RangeState) error {
	q := squirrel.Update(dmodels.RangeStatesTable).
		Where(squirrel.Eq{"rst_title": item.Title}).
		SetMap(map[string]interface{}{
			"rst_value_1d":   item.Value1d,
			"rst_value_7d":   item.Value7d,
			"rst_value_30d":  item.Value30d,
			"rst_value_90d":  item.Value90d,
			"rst_updated_at": item.UpdatedAt,
		})
	return m.update(q)
}

func (m DB) CreateRangeState(item dmodels.RangeState) error {
	q := squirrel.Insert(dmodels.RangeStatesTable).
		SetMap(map[string]interface{}{
			"rst_title":      item.Title,
			"rst_value_1d":   item.Value1d,
			"rst_value_7d":   item.Value7d,
			"rst_value_30d":  item.Value30d,
			"rst_value_90d":  item.Value90d,
			"rst_updated_at": item.UpdatedAt,
		})
	_, err := m.insert(q)
	return err
}
