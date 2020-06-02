package clickhouse

import (
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (db DB) CreateHistoricalStates(states []dmodels.HistoricalState) error {
	if len(states) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.HistoricalStates).Columns(
		"his_price",
		"his_market_cap",
		"his_circulation_supply",
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
			state.CirculationSupply,
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
