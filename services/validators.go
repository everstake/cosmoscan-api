package services

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/everstake/cosmoscan-api/smodels"
	"sort"
	"time"
)

const validatorsMapCacheKey = "validators_map"
const validatorsCacheKey = "validators"

func (s *ServiceFacade) UpdateValidatorsMap() {
	mp, err := s.makeValidatorMap()
	if err != nil {
		log.Error("UpdateValidatorsMap: makeValidatorMap: %s", err.Error())
		return
	}
	s.dao.CacheSet(validatorsMapCacheKey, mp, time.Minute*30)
}

func (s *ServiceFacade) GetValidatorMap() (map[string]node.Validator, error) {
	data, found := s.dao.CacheGet(validatorsMapCacheKey)
	if found {
		return data.(map[string]node.Validator), nil
	}
	mp, err := s.makeValidatorMap()
	if err != nil {
		return nil, fmt.Errorf("makeValidatorMap: %s", err.Error())
	}
	return mp, nil
}

func (s *ServiceFacade) makeValidatorMap() (map[string]node.Validator, error) {
	mp := make(map[string]node.Validator)
	validators, err := s.node.GetValidators()
	if err != nil {
		return nil, fmt.Errorf("node.GetValidators: %s", err.Error())
	}
	for _, validator := range validators {
		mp[validator.OperatorAddress] = validator
	}
	return mp, nil
}

func (s *ServiceFacade) GetStakingPie() (pie smodels.Pie, err error) {
	stakingPool, err := s.node.GetStakingPool()
	if err != nil {
		return pie, fmt.Errorf("node.GetStakingPool: %s", err.Error())
	}
	pie.Total = stakingPool.Result.BondedTokens
	validatorsMap, err := s.GetValidatorMap()
	if err != nil {
		return pie, fmt.Errorf("s.GetValidatorMap: %s", err.Error())
	}
	var validators []node.Validator
	for _, v := range validatorsMap {
		validators = append(validators, v)
	}
	sort.Slice(validators, func(i, j int) bool {
		return validators[i].DelegatorShares.GreaterThan(validators[j].DelegatorShares)
	})
	if len(validators) < 20 {
		return pie, fmt.Errorf("not enought validators")
	}
	parts := make([]smodels.PiePart, 20)
	for i := 0; i < 20; i++ {
		parts[i] = smodels.PiePart{
			Label: validators[i].OperatorAddress,
			Title: validators[i].Description.Moniker,
			Value: validators[i].DelegatorShares.Div(node.PrecisionDiv),
		}
	}
	pie.Parts = parts
	return pie, nil
}

func (s *ServiceFacade) GetValidators() (validators []smodels.Validator, err error) {
	data, found := s.dao.CacheGet(validatorsCacheKey)
	if found {
		return data.([]smodels.Validator), nil
	}
	return nil, fmt.Errorf("not found in cache")
}

func (s *ServiceFacade) UpdateValidators() {
	validators, err := s.makeValidators()
	if err != nil {
		log.Error("UpdateValidators: makeValidators: %s", err.Error())
		return
	}
	s.dao.CacheSet(validatorsCacheKey, validators, time.Hour)
}

func (s *ServiceFacade) makeValidators() (validators []smodels.Validator, err error) {
	nodeValidators, err := s.node.GetValidators()
	if err != nil {
		return nil, fmt.Errorf("node.GetValidators: %s", err.Error())
	}
	for _, v := range nodeValidators {
		key, err := types.GetPubKeyFromBech32(types.Bech32PubKeyTypeConsPub, v.ConsensusPubkey)
		if err != nil {
			return nil, fmt.Errorf("types.GetPubKeyFromBech32: %s", err.Error())
		}

		blockProposed, err := s.dao.GetProposedBlocksTotal(filters.BlocksProposed{Proposers: []string{key.Address().String()}})
		if err != nil {
			return nil, fmt.Errorf("dao.GetProposedBlocksTotal: %s", err.Error())
		}

		addressBytes, err := types.GetFromBech32(v.OperatorAddress, types.Bech32PrefixValAddr)
		if err != nil {
			return nil, fmt.Errorf("types.GetFromBech32: %s", err.Error())
		}
		address, err := types.AccAddressFromHex(hex.EncodeToString(addressBytes))
		if err != nil {
			return nil, fmt.Errorf("types.AccAddressFromHex: %s", err.Error())
		}

		totalVotes, err := s.dao.GetProposalVotesTotal(filters.ProposalVotes{
			Voters: []string{address.String()},
		})
		if err != nil {
			return nil, fmt.Errorf("dao.GetProposalVotesTotal: %s", err.Error())
		}

		delegatorsTotal, err := s.dao.GetDelegatorsTotal(filters.Delegators{Validators: []string{v.OperatorAddress}})
		if err != nil {
			return nil, fmt.Errorf("dao.GetDelegatorsTotal: %s", err.Error())
		}

		power24Change, err := s.dao.GetVotingPower(filters.VotingPower{
			TimeRange: filters.TimeRange{
				From: dmodels.NewTime(time.Now().Add(-time.Hour * 24)),
				To:   dmodels.NewTime(time.Now()),
			},
			Validators: []string{v.OperatorAddress},
		})

		selfStake, err := s.node.GetDelegatorValidatorStake(address.String(), v.OperatorAddress)
		if err != nil {
			return nil, fmt.Errorf("node.GetDelegatorValidatorStake: %s", err.Error())
		}

		validators = append(validators, smodels.Validator{
			Title:           v.Description.Moniker,
			Power:           v.DelegatorShares.Div(node.PrecisionDiv),
			SelfStake:       selfStake,
			Fee:             v.Commission.CommissionRates.Rate,
			BlocksProposed:  blockProposed,
			Delegators:      delegatorsTotal,
			Power24Change:   power24Change,
			GovernanceVotes: totalVotes,
		})
	}

	return validators, nil
}
