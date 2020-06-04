package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (db DB) CreateProposals(proposals []dmodels.Proposal) error {
	if len(proposals) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.ProposalsTable).Columns(
		"pro_id",
		"pro_title",
		"pro_description",
		"pro_recipient",
		"pro_amount",
		"pro_init_deposit",
		"pro_proposer",
		"pro_created_at",
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

func (db DB) GetProposals(filter filters.Proposals) (proposals []dmodels.Proposal, err error) {
	q := squirrel.Select("*").From(dmodels.ProposalsTable).OrderBy("pro_created_at desc")
	if filter.Limit != 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset != 0 {
		q = q.Limit(filter.Offset)
	}
	err = db.Find(&proposals, q)
	return proposals, err
}
