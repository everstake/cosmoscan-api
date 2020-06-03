package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/smodels"
)

func (db DB) CreateHistoricalStates(states []dmodels.HistoricalState) error {
	if len(states) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.HistoricalStates).Columns(
		"his_price",
		"his_market_cap",
		"his_circulating_supply",
		"his_trading_volume",
		"his_staked_ratio",
		"his_inflation_rate",
		"his_transactions_count",
		"his_community_pool",
		"his_top_20_weight",
		"his_created_at",
	)
	for _, state := range states {
		q = q.Values(
			state.Price,
			state.MarketCap,
			state.CirculatingSupply,
			state.TradingVolume,
			state.StakedRatio,
			state.InflationRate,
			state.TransactionsCount,
			state.CommunityPool,
			state.Top20Weight,
			state.CreatedAt,
		)
	}
	return db.Insert(q)
}

func (db DB) GetHistoricalStates(state filters.HistoricalState) (states []dmodels.HistoricalState, err error) {
	q := squirrel.Select("*").From(dmodels.HistoricalStates).OrderBy("his_created_at desc")
	if state.Limit != 0 {
		q = q.Limit(state.Limit)
	}
	if state.Offset != 0 {
		q = q.Limit(state.Offset)
	}
	err = db.Find(&states, q)
	return states, err
}

func (db DB) GetAggHistoricalStatesByField(filter filters.Agg, field string) (items []smodels.AggItem, err error) {
	q := squirrel.Select(
		fmt.Sprintf("avg(%s) AS value", field),
		fmt.Sprintf("toDateTime(%s(his_created_at)) AS time", filter.AggFunc()),
	).From(dmodels.HistoricalStates).
		GroupBy("time").
		OrderBy("time")
	if !filter.From.IsZero() {
		q = q.Where(squirrel.GtOrEq{"his_created_at": filter.From.Time})
	}
	if !filter.To.IsZero() {
		q = q.Where(squirrel.LtOrEq{"his_created_at": filter.To.Time})
	}
	err = db.Find(&items, q)
	if err != nil {
		return nil, err
	}
	return items, nil
}
