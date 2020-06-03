package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
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
	return blocks, nil
}
