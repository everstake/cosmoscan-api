package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (db DB) CreateMissedBlocks(blocks []dmodels.MissedBlock) error {
	if len(blocks) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.MissedBlocks).Columns("mib_id", "mib_height", "mib_validator", "mib_created_at")
	for _, block := range blocks {
		if block.ID == "" {
			return fmt.Errorf("field ProposalID can not be empty")
		}
		if block.Height == 0 {
			return fmt.Errorf("field ProposalID can not be zero")
		}
		if block.Validator == "" {
			return fmt.Errorf("field Validator can not be empty")
		}
		if block.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be 0")
		}
		q = q.Values(block.ID, block.Height, block.Validator, block.CreatedAt)
	}
	return db.Insert(q)
}

func (db DB) GetMissedBlocksCount(filter filters.MissedBlocks) (total uint64, err error) {
	q := squirrel.Select("count(*) as total").From(dmodels.MissedBlocks)
	if len(filter.Validators) != 0 {
		q = q.Where(squirrel.Eq{"mib_validator": filter.Validators})
	}
	err = db.FindFirst(&total, q)
	return total, err
}
