package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (db DB) CreateHistoryProposals(proposals []dmodels.HistoryProposal) error {
	if len(proposals) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.HistoryProposalsTable).Columns(
		"hpr_id",
		"hpr_tx_hash",
		"hpr_title",
		"hpr_description",
		"hpr_recipient",
		"hpr_amount",
		"hpr_init_deposit",
		"hpr_proposer",
		"hpr_created_at",
	)
	for _, proposal := range proposals {
		if proposal.ID == 0 {
			return fmt.Errorf("field ID can not be 0")
		}
		if proposal.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(
			proposal.ID,
			proposal.TxHash,
			proposal.Title,
			proposal.Description,
			proposal.Recipient,
			proposal.Amount,
			proposal.InitDeposit,
			proposal.Proposer,
			proposal.CreatedAt,
		)
	}
	return db.Insert(q)
}

func (db DB) GetHistoryProposals(filter filters.HistoryProposals) (proposals []dmodels.HistoryProposal, err error) {
	q := squirrel.Select("*").From(dmodels.HistoryProposalsTable).OrderBy("hpr_created_at desc")
	if len(filter.ID) != 0 {
		q = q.Where(squirrel.Eq{"hpr_id": filter.ID})
	}
	if filter.Limit != 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset != 0 {
		q = q.Limit(filter.Offset)
	}
	err = db.Find(&proposals, q)
	return proposals, err
}
