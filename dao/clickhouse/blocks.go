package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/smodels"
)

func (db DB) CreateBlocks(blocks []dmodels.Block) error {
	if len(blocks) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.BlocksTable).Columns("blk_id", "blk_hash", "blk_proposer", "blk_created_at")
	for _, block := range blocks {
		if block.ID == 0 {
			return fmt.Errorf("field ID can not be 0")
		}
		if block.Hash == "" {
			return fmt.Errorf("hash can not be empty")
		}
		if block.Proposer == "" {
			return fmt.Errorf("proposer can not be empty")
		}
		if block.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be 0")
		}
		q = q.Values(block.ID, block.Hash, block.Proposer, block.CreatedAt)
	}
	return db.Insert(q)
}

func (db DB) GetBlocks(filter filters.Blocks) (blocks []dmodels.Block, err error) {
	q := squirrel.Select("*").From(dmodels.BlocksTable).OrderBy("blk_id desc")
	if filter.Limit != 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset != 0 {
		q = q.Offset(filter.Offset)
	}
	err = db.Find(&blocks, q)
	return blocks, err
}

func (db DB) GetAggBlocksCount(filter filters.Agg) (items []smodels.AggItem, err error) {
	q := filter.BuildQuery("toDecimal64(count(blk_id), 0)", "blk_created_at", dmodels.BlocksTable)
	err = db.Find(&items, q)
	return items, err
}

func (db DB) GetAggBlocksDelay(filter filters.Agg) (items []smodels.AggItem, err error) {
	q := squirrel.Select(
		"avg(toUnixTimestamp(b1.blk_created_at) - toUnixTimestamp(b2.blk_created_at)) as value",
		fmt.Sprintf("toDateTime(%s(b1.blk_created_at)) AS time", filter.AggFunc()),
	).From(fmt.Sprintf("%s as b1", dmodels.BlocksTable)).
		JoinClause("JOIN blocks as b2 ON b1.blk_id = toUInt64(plus(b2.blk_id, 1))").
		Where(squirrel.Gt{"b1.blk_id": 2}).
		GroupBy("time").
		OrderBy("time")

	if !filter.From.IsZero() {
		q = q.Where(squirrel.GtOrEq{"time": filter.From.Time})
	}
	if !filter.To.IsZero() {
		q = q.Where(squirrel.LtOrEq{"time": filter.To.Time})
	}
	err = db.Find(&items, q)
	return items, err
}

func (db DB) GetAggUniqBlockValidators(filter filters.Agg) (items []smodels.AggItem, err error) {
	q := filter.BuildQuery("toDecimal64(count(DISTINCT blk_proposer), 0)", "blk_created_at", dmodels.BlocksTable)
	err = db.Find(&items, q)
	return items, err
}

func (db DB) GetAvgBlocksDelay(filter filters.TimeRange) (delay float64, err error) {
	q := squirrel.Select(
		"avg(toUnixTimestamp(b1.blk_created_at) - toUnixTimestamp(b2.blk_created_at)) as delay",
	).From(fmt.Sprintf("%s as b1", dmodels.BlocksTable)).
		JoinClause("JOIN blocks as b2 ON b1.blk_id = toUInt64(plus(b2.blk_id, 1))").
		Where(squirrel.Gt{"b1.blk_id": 2})
	if !filter.From.IsZero() {
		q = q.Where(squirrel.GtOrEq{"b1.blk_created_at": filter.From.Time})
	}
	if !filter.To.IsZero() {
		q = q.Where(squirrel.LtOrEq{"b1.blk_created_at": filter.To.Time})
	}
	err = db.FindFirst(&delay, q)
	return delay, err
}
