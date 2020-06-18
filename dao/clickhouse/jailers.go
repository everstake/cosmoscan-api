package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (db DB) CreateJailers(jailers []dmodels.Jailer) error {
	if len(jailers) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.JailersTable).Columns("jlr_id", "jlr_address", "jlr_created_at")
	for _, jailer := range jailers {
		if jailer.ID == "" {
			return fmt.Errorf("field ProposalID can not be empty")
		}
		if jailer.Address == "" {
			return fmt.Errorf("field Address can not be empty")
		}
		if jailer.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(jailer.ID, jailer.Address, jailer.CreatedAt)
	}
	return db.Insert(q)
}

func (db DB) GetJailersTotal() (total uint64, err error) {
	q := squirrel.Select("count(*) as total").From(dmodels.JailersTable)
	err = db.FindFirst(&total, q)
	return total, err
}
