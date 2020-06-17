package mysql

import (
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
)

func (m DB) CreateProposals(proposals []dmodels.Proposal) error {
	if len(proposals) == 0 {
		return nil
	}
	q := squirrel.Insert(dmodels.ProposalsTable).Columns(
		"pro_id",
		"pro_proposer",
		"pro_title",
		"pro_description",
		"pro_status",
		"pro_votes_yes",
		"pro_votes_abstain",
		"pro_votes_no",
		"pro_votes_no_with_veto",
		"pro_submit_time",
		"pro_deposit_end_time",
		"pro_total_deposits",
		"pro_voting_start_time",
		"pro_voting_end_time",
	)
	for _, p := range proposals {
		if p.ID == 0 {
			return fmt.Errorf("invalid ID")
		}

		q = q.Values(
			p.ID,
			p.Proposer,
			p.Title,
			p.Description,
			p.Status,
			p.VotesYes,
			p.VotesAbstain,
			p.VotesNo,
			p.VotesNoWithVeto,
			p.SubmitTime,
			p.DepositEndTime,
			p.TotalDeposits,
			p.VotingStartTime,
			p.VotingEndTime,
		)
	}
	_, err := m.insert(q)
	return err
}

func (m DB) GetProposals(filter filters.Proposals) (proposals []dmodels.Proposal, err error) {
	q := squirrel.Select("*").From(dmodels.ProposalsTable).OrderBy("pro_id desc")
	if len(filter.ID) != 0 {
		q = q.Where(squirrel.Eq{"pro_id": filter.ID})
	}
	if filter.Limit != 0 {
		q = q.Limit(filter.Limit)
	}
	if filter.Offset != 0 {
		q = q.Limit(filter.Offset)
	}
	err = m.find(&proposals, q)
	return proposals, err
}

func (m DB) UpdateProposal(proposal dmodels.Proposal) error {
	q := squirrel.Update(dmodels.ProposalsTable).
		Where(squirrel.Eq{"pro_id": proposal.ID}).
		SetMap(map[string]interface{}{
			"pro_proposer":           proposal.Proposer,
			"pro_title":              proposal.Title,
			"pro_description":        proposal.Description,
			"pro_status":             proposal.Status,
			"pro_votes_yes":          proposal.VotesYes,
			"pro_votes_abstain":      proposal.VotesAbstain,
			"pro_votes_no":           proposal.VotesNo,
			"pro_votes_no_with_veto": proposal.VotesNoWithVeto,
			"pro_submit_time":        proposal.SubmitTime,
			"pro_deposit_end_time":   proposal.DepositEndTime,
			"pro_total_deposits":     proposal.TotalDeposits,
			"pro_voting_start_time":  proposal.VotingStartTime,
			"pro_voting_end_time":    proposal.VotingEndTime,
		})
	return m.update(q)
}
