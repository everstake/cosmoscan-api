package clickhouse

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/smodels"
)

func (db DB) CreateProposalVotes(votes []dmodels.ProposalVote) error {
	if len(votes) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.ProposalVotesTable).Columns("prv_id", "prv_proposal_id", "prv_voter", "prv_tx_hash", "prv_option", "prv_created_at")
	for _, vote := range votes {
		if vote.ID == "" {
			return fmt.Errorf("field ProposalID can not be empty")
		}
		if vote.ProposalID == 0 {
			return fmt.Errorf("field ProposalID can not be zero")
		}
		if vote.Voter == "" {
			return fmt.Errorf("field Voter can not be empty")
		}
		if vote.TxHash == "" {
			return fmt.Errorf("field TxHash can not be empty")
		}
		if vote.CreatedAt.IsZero() {
			return fmt.Errorf("field CreatedAt can not be zero")
		}
		q = q.Values(vote.ID, vote.ProposalID, vote.Voter, vote.TxHash, vote.Option, vote.CreatedAt)
	}
	return db.Insert(q)
}

func (db DB) GetProposalVotes(filter filters.ProposalVotes) (votes []dmodels.ProposalVote, err error) {
	q := squirrel.Select("*").From(dmodels.ProposalVotesTable).OrderBy("prv_created_at")
	if filter.ProposalID != 0 {
		q = q.Where(squirrel.Eq{"prv_proposal_id": filter.ProposalID})
	}
	if len(filter.Voters) != 0 {
		q = q.Where(squirrel.Eq{"prv_voter": filter.Voters})
	}
	if filter.Limit != 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset != 0 {
		q = q.Offset(filter.Offset)
	}
	err = db.Find(&votes, q)
	return votes, err
}

func (db DB) GetAggProposalVotes(filter filters.Agg, id []uint64) (items []smodels.AggItem, err error) {
	q := filter.BuildQuery("toDecimal64(count(*), 0)", "prv_created_at", dmodels.ProposalVotesTable)
	if len(id) != 0 {
		q = q.Where(squirrel.Eq{"prv_proposal_id": id})
	}
	err = db.Find(&items, q)
	return items, err
}

func (db DB) GetTotalVotesByAddress(address string) (total uint64, err error) {
	q := squirrel.Select("count(distinct prv_proposal_id) as total").
		From(dmodels.ProposalVotesTable).
		Where(squirrel.Eq{"prv_voter": address})
	err = db.FindFirst(&total, q)
	return total, err
}
