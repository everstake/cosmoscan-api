package services

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
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

	totalAccounts, err := s.dao.GetAccountsTotal(filters.Accounts{})
	if err != nil {
		log.Error("UpdateProposals: node.GetProposals: %s", err.Error())
		return
	}

	validators, err := s.node.GetValidators()
	if err != nil {
		log.Error("UpdateProposals: node.GetValidators: %s", err.Error())
		return
	}
	validatorsMap := make(map[string]node.Validator)
	for _, validator := range validators {
		bench, _ := types.ValAddressFromBech32(validator.OperatorAddress)
		accAddress := types.AccAddress(bench.Bytes())
		validatorsMap[accAddress.String()] = validator
	}

	totalStake, err := s.node.GetStakingPool()
	if err != nil {
		log.Error("UpdateProposals: node.GetStakingPool: %s", err.Error())
		return
	}

	for _, p := range nodeProposals.Result {
		votersTotal, err := s.dao.GetProposalVotesTotal(filters.ProposalVotes{ProposalID: []uint64{p.ID}})
		if err != nil {
			log.Error("UpdateProposals: node.GetProposalVotesTotal: %s", err.Error())
			return
		}
		participationRate := decimal.Zero
		if votersTotal != 0 {
			participationRate = decimal.NewFromFloat(float64(votersTotal) / float64(totalAccounts) * 100).Truncate(2)
		}

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

		activityItems, err := s.dao.GetAggProposalVotes(filters.Agg{
			By: filters.AggByDay,
		}, []uint64{p.ID})
		activityJson, _ := json.Marshal(activityItems)

		hps, err := s.dao.GetHistoryProposals(filters.HistoryProposals{ID: []uint64{p.ID}})
		if err != nil {
			log.Error("UpdateProposals: dao.GetHistoryProposals: %s", err.Error())
			return
		}
		var txHash string
		if len(hps) > 0 {
			txHash = hps[0].TxHash
		}

		yes := decimal.NewFromInt(p.FinalTallyResult.Yes).Div(node.PrecisionDiv)
		abstain := decimal.NewFromInt(p.FinalTallyResult.Abstain).Div(node.PrecisionDiv)
		no := decimal.NewFromInt(p.FinalTallyResult.No).Div(node.PrecisionDiv)
		noWithVeto := decimal.NewFromInt(p.FinalTallyResult.NoWithVeto).Div(node.PrecisionDiv)

		turnout := decimal.Zero
		if !totalStake.Result.BondedTokens.IsZero() {
			turnout = yes.Add(abstain).Add(no).Add(noWithVeto).Div(totalStake.Result.BondedTokens)
			turnout = turnout.Mul(decimal.NewFromFloat(100)).Truncate(2)
		}

		proposal := dmodels.Proposal{
			ID:                p.ID,
			TxHash:            txHash,
			Proposer:          proposer,
			Type:              p.Content.Type,
			Title:             p.Content.Value.Title,
			Description:       p.Content.Value.Description,
			Status:            p.ProposalStatus,
			VotesYes:          yes,
			VotesAbstain:      abstain,
			VotesNo:           no,
			VotesNoWithVeto:   noWithVeto,
			SubmitTime:        dmodels.NewTime(p.SubmitTime),
			DepositEndTime:    dmodels.NewTime(p.DepositEndTime),
			TotalDeposits:     totalDeposit.Div(node.PrecisionDiv),
			VotingStartTime:   dmodels.NewTime(p.VotingStartTime),
			VotingEndTime:     dmodels.NewTime(p.VotingEndTime),
			Voters:            votersTotal,
			ParticipationRate: participationRate,
			Turnout:           turnout,
			Activity:          activityJson,
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
