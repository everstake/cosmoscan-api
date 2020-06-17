package services

import (
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/shopspring/decimal"
)

func (s *ServiceFacade) GetProposals(filter filters.Proposals) (proposals []dmodels.Proposal, err error) {
	proposals, err = s.dao.GetProposals(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetProposals: %s", err.Error())
	}
	return proposals, nil
}

func (s *ServiceFacade) UpdateProposals() {
	nodeProposals, err := s.node.GetProposals()
	if err != nil {
		log.Error("UpdateProposals: node.GetProposals: %s", err.Error())
		return
	}
	for _, p := range nodeProposals.Result {
		proposer, err := s.node.GetProposalProposer(p.ID)
		if err != nil {
			log.Error("UpdateProposals: node.GetProposalProposer: %s", err.Error())
			return
		}
		proposals, err := s.dao.GetProposals(filters.Proposals{
			ID:    []uint64{p.ID},
			Limit: 1,
		})
		if err != nil {
			log.Error("UpdateProposals: dao.GetProposals: %s", err.Error())
			return
		}
		totalDeposit := decimal.Zero
		for _, value := range p.TotalDeposit {
			totalDeposit = totalDeposit.Add(value.Amount)
		}
		proposal := dmodels.Proposal{
			ID:              p.ID,
			Proposer:        proposer,
			Title:           p.Content.Value.Title,
			Description:     p.Content.Value.Description,
			Status:          p.ProposalStatus,
			VotesYes:        decimal.NewFromInt(p.FinalTallyResult.Yes).Div(node.PrecisionDiv),
			VotesAbstain:    decimal.NewFromInt(p.FinalTallyResult.Abstain).Div(node.PrecisionDiv),
			VotesNo:         decimal.NewFromInt(p.FinalTallyResult.No).Div(node.PrecisionDiv),
			VotesNoWithVeto: decimal.NewFromInt(p.FinalTallyResult.NoWithVeto).Div(node.PrecisionDiv),
			SubmitTime:      p.SubmitTime,
			DepositEndTime:  p.DepositEndTime,
			TotalDeposits:   totalDeposit.Div(node.PrecisionDiv),
			VotingStartTime: p.VotingStartTime,
			VotingEndTime:   p.VotingEndTime,
		}
		if len(proposals) == 0 {
			err = s.dao.CreateProposals([]dmodels.Proposal{proposal})
		} else {
			err = s.dao.UpdateProposal(proposal)
		}
		if err != nil {
			log.Error("UpdateProposals: save/update proposal: %s", err.Error())
			return
		}
	}
}
