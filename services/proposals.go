package services

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
	"time"
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
		votes, err := s.GetProposalVotes(filters.ProposalVotes{ProposalID: p.ID})
		if err != nil {
			log.Error("UpdateProposals: GetProposalVotes: %s", err.Error())
			return
		}
		votersTotal := len(votes)
		participationRate := decimal.Zero
		if votersTotal != 0 {
			participationRate = decimal.NewFromFloat(float64(votersTotal) / float64(totalAccounts) * 100).Truncate(2)
		}

		proposerAddress, err := s.node.GetProposalProposer(p.ID)
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

		var yes, abstain, no, noWithVeto decimal.Decimal
		if p.ProposalStatus == "VotingPeriod" {
			tally, err := s.node.ProposalTallyResult(p.ID)
			if err != nil {
				log.Error("UpdateProposals: node.ProposalTallyResult: %s", err.Error())
				return
			}
			yes = decimal.NewFromInt(tally.Result.Yes).Div(node.PrecisionDiv)
			abstain = decimal.NewFromInt(tally.Result.Abstain).Div(node.PrecisionDiv)
			no = decimal.NewFromInt(tally.Result.No).Div(node.PrecisionDiv)
			noWithVeto = decimal.NewFromInt(tally.Result.NoWithVeto).Div(node.PrecisionDiv)
		} else {
			yes = decimal.NewFromInt(p.FinalTallyResult.Yes).Div(node.PrecisionDiv)
			abstain = decimal.NewFromInt(p.FinalTallyResult.Abstain).Div(node.PrecisionDiv)
			no = decimal.NewFromInt(p.FinalTallyResult.No).Div(node.PrecisionDiv)
			noWithVeto = decimal.NewFromInt(p.FinalTallyResult.NoWithVeto).Div(node.PrecisionDiv)
		}

		turnout := decimal.Zero
		if !totalStake.Result.BondedTokens.IsZero() {
			turnout = yes.Add(abstain).Add(no).Add(noWithVeto).Div(totalStake.Result.BondedTokens)
			turnout = turnout.Mul(decimal.NewFromFloat(100)).Truncate(2)
		}

		proposer := proposerAddress
		a, ok := validatorsMap[proposerAddress]
		if ok {
			proposer = a.Description.Moniker
			proposerAddress = a.OperatorAddress
		}

		proposal := dmodels.Proposal{
			ID:                p.ID,
			TxHash:            txHash,
			Type:              p.Content.Type,
			Proposer:          proposer,
			ProposerAddress:   proposerAddress,
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
			Voters:            uint64(votersTotal),
			ParticipationRate: participationRate,
			Turnout:           turnout,
			Activity:          activityJson,
		}

		if proposal.VotingStartTime.Unix() < 0 {
			proposal.VotingStartTime = dmodels.Time{Time: time.Unix(0, 0)}
		}
		if proposal.VotingEndTime.Unix() < 0 {
			proposal.VotingEndTime = dmodels.Time{Time: time.Unix(0, 0)}
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

func (s *ServiceFacade) GetProposalVotes(filter filters.ProposalVotes) (items []smodels.ProposalVote, err error) {
	votes, err := s.dao.GetProposalVotes(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetProposalVotes: %s", err.Error())
	}
	vm, err := s.GetValidatorMap()
	if err != nil {
		return nil, fmt.Errorf("GetValidatorMap: %s", err.Error())
	}
	validatorsMap := make(map[string]node.Validator)
	for _, validator := range vm {
		bench, _ := types.ValAddressFromBech32(validator.OperatorAddress)
		accAddress := types.AccAddress(bench.Bytes())
		validatorsMap[accAddress.String()] = validator
	}
	votesMap := make(map[string]dmodels.ProposalVote)
	for _, vote := range votes {
		votesMap[vote.Voter] = vote
	}
	for _, vote := range votesMap {
		title := vote.Voter
		var isValidator bool
		validator, ok := validatorsMap[vote.Voter]
		if ok {
			title = validator.Description.Moniker
			isValidator = ok
		}
		items = append(items, smodels.ProposalVote{
			Title:        title,
			IsValidator:  isValidator,
			ProposalVote: vote,
		})
	}
	return items, nil
}

func (s *ServiceFacade) GetProposalDeposits(filter filters.ProposalDeposits) (deposits []dmodels.ProposalDeposit, err error) {
	deposits, err = s.dao.GetProposalDeposits(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetProposalDeposits: %s", err.Error())
	}
	return deposits, nil
}

func (s *ServiceFacade) GetProposalsChartData() (items []smodels.ProposalChartData, err error) {
	proposals, err := s.dao.GetProposals(filters.Proposals{})
	if err != nil {
		return nil, fmt.Errorf("dao.GetProposals: %s", err.Error())
	}
	validators, err := s.GetValidatorMap()
	if err != nil {
		return nil, fmt.Errorf("GetValidatorMap: %s", err.Error())
	}
	validatorsMap := make(map[string]node.Validator)
	for _, validator := range validators {
		bench, _ := types.ValAddressFromBech32(validator.OperatorAddress)
		accAddress := types.AccAddress(bench.Bytes())
		validatorsMap[accAddress.String()] = validator
	}

	for _, p := range proposals {
		votes, err := s.dao.GetProposalVotes(filters.ProposalVotes{ProposalID: p.ID})
		if err != nil {
			return nil, fmt.Errorf("dao.GetProposalVotes: %s", err.Error())
		}
		var validatorsTotal uint64
		for _, vote := range votes {
			_, ok := validatorsMap[vote.Voter]
			if ok {
				validatorsTotal++
			}
		}

		totalAmount := p.VotesYes.Add(p.VotesNo).Add(p.VotesAbstain).Add(p.VotesNoWithVeto)

		pd := smodels.ProposalChartData{
			ProposalID:      p.ID,
			VotersTotal:     uint64(len(votes)),
			ValidatorsTotal: validatorsTotal,
			Turnout:         p.Turnout,
		}

		if !totalAmount.IsZero() {
			d100 := decimal.New(100, 0)
			pd.YesPercent = p.VotesYes.Div(totalAmount).Mul(d100).Truncate(2)
			pd.NoPercent = p.VotesNo.Div(totalAmount).Mul(d100).Truncate(2)
			pd.NoWithVetoPercent = p.VotesNoWithVeto.Div(totalAmount).Mul(d100).Truncate(2)
			pd.AbstainPercent = p.VotesAbstain.Div(totalAmount).Mul(d100).Truncate(2)
		}

		items = append(items, pd)
	}

	return items, nil
}
