package services

import (
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services/helpers"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
	"sort"
	"time"
)

const validatorsMapCacheKey = "validators_map"
const validatorsCacheKey = "validators"
const mostJailedValidators = "mostJailedValidators"
const Bech32PrefixValAddr = "bitsongvaloper"

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

func (s *ServiceFacade) getConsensusValidatorMap() (map[string]node.Validator, error) {
	vMap, err := s.GetValidatorMap()
	if err != nil {
		return nil, fmt.Errorf("GetValidatorMap: %s", err.Error())
	}
	cvMap := make(map[string]node.Validator)
	for _, v := range vMap {
		consAddress, err := helpers.GetHexAddressFromBase64PK(v.ConsensusPubkey.Key)
		if err != nil {
			return nil, fmt.Errorf("helpers.GetHexAddressFromBase64PK: %s", err.Error())
		}
		cvMap[consAddress] = v
	}
	return cvMap, nil
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
	pie.Total = stakingPool.Pool.BondedTokens
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
	stakingPool, err := s.node.GetStakingPool()
	if err != nil {
		return nil, fmt.Errorf("node.GetStakingPool: %s", err.Error())
	}
	for _, v := range nodeValidators {
		consAddress, err := helpers.GetHexAddressFromBase64PK(v.ConsensusPubkey.Key)
		if err != nil {
			return nil, fmt.Errorf("helpers.GetHexAddressFromBase64PK: %s", err.Error())
		}
		blockProposed, err := s.dao.GetProposedBlocksTotal(filters.BlocksProposed{Proposers: []string{consAddress}})
		if err != nil {
			return nil, fmt.Errorf("dao.GetProposedBlocksTotal: %s", err.Error())
		}

		addressBytes, err := types.GetFromBech32(v.OperatorAddress, Bech32PrefixValAddr)
		if err != nil {
			return nil, fmt.Errorf("types.GetFromBech32: %s", err.Error())
		}
		address, err := types.AccAddressFromHex(hex.EncodeToString(addressBytes))
		if err != nil {
			return nil, fmt.Errorf("types.AccAddressFromHex: %s", err.Error())
		}
		totalVotes, err := s.dao.GetTotalVotesByAddress(address.String())
		if err != nil {
			return nil, fmt.Errorf("dao.GetTotalVotesByAddress: %s", err.Error())
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

		power := v.DelegatorShares.Div(node.PrecisionDiv)
		percentPower := decimal.Zero
		if !stakingPool.Pool.BondedTokens.IsZero() {
			percentPower = power.Div(stakingPool.Pool.BondedTokens).Mul(decimal.NewFromInt(100)).Truncate(2)
		}

		validators = append(validators, smodels.Validator{
			Title:           v.Description.Moniker,
			Power:           power,
			PercentPower:    percentPower,
			SelfStake:       selfStake,
			Fee:             v.Commission.CommissionRates.Rate,
			BlocksProposed:  blockProposed,
			Delegators:      delegatorsTotal,
			Power24Change:   power24Change,
			GovernanceVotes: totalVotes,
			Website:         v.Description.Website,
			OperatorAddress: v.OperatorAddress,
			AccAddress:      address.String(),
			ConsAddress:     consAddress,
		})
	}

	return validators, nil
}

func (s *ServiceFacade) GetTopProposedBlocksValidators() (items []dmodels.ValidatorValue, err error) {
	data, found := s.dao.CacheGet(topProposedBlocksValidatorsKey)
	if found {
		return data.([]dmodels.ValidatorValue), nil
	}
	items, err = s.dao.GetTopProposedBlocksValidators()
	if err != nil {
		return nil, fmt.Errorf("dao.GetTopProposedBlocksValidators: %s", err.Error())
	}
	validators, err := s.GetValidatorMap()
	if err != nil {
		return nil, fmt.Errorf("GetValidators: %s", err.Error())
	}
	mp := make(map[string]string)
	for _, validator := range validators {
		address, err := helpers.GetHexAddressFromBase64PK(validator.ConsensusPubkey.Key)
		if err != nil {
			return nil, fmt.Errorf("helpers.GetHexAddressFromBase64PK: %s", err.Error())
		}
		mp[address] = validator.Description.Moniker
	}
	for i, item := range items {
		title, found := mp[item.Validator]
		if found {
			items[i] = dmodels.ValidatorValue{
				Validator: title,
				Value:     item.Value,
			}
		}
	}
	s.dao.CacheSet(topProposedBlocksValidatorsKey, items, time.Minute*60)
	return items, nil
}

func (s *ServiceFacade) GetMostJailedValidators() (items []dmodels.ValidatorValue, err error) {
	data, found := s.dao.CacheGet(mostJailedValidators)
	if found {
		return data.([]dmodels.ValidatorValue), nil
	}
	items, err = s.dao.GetMostJailedValidators()
	if err != nil {
		return nil, fmt.Errorf("dao.GetMostJailedValidators: %s", err.Error())
	}
	validators, err := s.GetValidatorMap()
	if err != nil {
		return nil, fmt.Errorf("GetValidators: %s", err.Error())
	}
	mp := make(map[string]string)
	for _, validator := range validators {
		mp[validator.OperatorAddress] = validator.Description.Moniker
	}
	for i, item := range items {
		title, found := mp[item.Validator]
		if found {
			items[i] = dmodels.ValidatorValue{
				Validator: title,
				Value:     item.Value,
			}
		}
	}
	s.dao.CacheSet(mostJailedValidators, items, time.Minute*60)
	return items, nil
}

func (s *ServiceFacade) GetFeeRanges() (items []smodels.FeeRange, err error) {
	point := int64(10)
	min := decimal.Zero
	max := decimal.Zero
	validatorsMap, err := s.GetValidatorMap()
	for _, validator := range validatorsMap {
		if min.IsZero() && max.IsZero() {
			min = validator.Commission.CommissionRates.Rate
			max = validator.Commission.CommissionRates.Rate
			continue
		}
		if validator.Commission.CommissionRates.Rate.LessThan(min) {
			min = validator.Commission.CommissionRates.Rate
		}
		if validator.Commission.CommissionRates.Rate.GreaterThan(max) {
			max = validator.Commission.CommissionRates.Rate
		}
	}
	step := max.Sub(min).Div(decimal.NewFromInt(point))
	for i := int64(1); i <= point; i++ {
		var validators []smodels.FeeRangeValidator
		from := step.Mul(decimal.NewFromInt(i)).Sub(step)
		to := step.Mul(decimal.NewFromInt(i))

		for _, validator := range validatorsMap {
			rate := validator.Commission.CommissionRates.Rate
			if rate.GreaterThan(from) && rate.LessThanOrEqual(to) {
				validators = append(validators, smodels.FeeRangeValidator{
					Validator: validator.Description.Moniker,
					Fee:       rate,
				})
			}
		}
		items = append(items, smodels.FeeRange{
			From:       from,
			To:         to,
			Validators: validators,
		})
	}
	return items, nil
}

func (s *ServiceFacade) GetValidatorsDelegatorsTotal() (values []dmodels.ValidatorValue, err error) {
	validatorsMap, err := s.GetValidatorMap()
	if err != nil {
		return nil, fmt.Errorf("GetValidatorMap: %s", err.Error())
	}
	values, err = s.dao.GetValidatorsDelegatorsTotal()
	if err != nil {
		return nil, fmt.Errorf("dao.GetValidatorsDelegatorsTotal: %s", err.Error())
	}
	for i, v := range values {
		validator, found := validatorsMap[v.Validator]
		if found {
			values[i].Validator = validator.Description.Moniker
		}
	}
	return values, nil
}

func (s *ServiceFacade) GetValidator(address string) (validator smodels.Validator, err error) {
	data, found := s.dao.CacheGet(validatorsCacheKey)
	if !found {
		return validator, fmt.Errorf("not found in cache")
	}
	validators, ok := data.([]smodels.Validator)
	if !ok {
		return validator, fmt.Errorf("can`t cast to current type")
	}
	for _, v := range validators {
		if v.OperatorAddress == address {
			return v, nil
		}
	}
	return validator, fmt.Errorf("not found validator with address: %s", address)
}

func (s *ServiceFacade) GetValidatorBalance(valAddress string) (balance smodels.Balance, err error) {
	validator, err := s.GetValidator(valAddress)
	if err != nil {
		return balance, fmt.Errorf("GetValidator: %s", err.Error())
	}
	balance.SelfDelegated = validator.SelfStake
	balance.OtherDelegated = validator.Power.Sub(validator.SelfStake)
	addressBytes, err := types.GetFromBech32(valAddress, Bech32PrefixValAddr)
	if err != nil {
		return balance, fmt.Errorf("types.GetFromBech32: %s", err.Error())
	}
	address, err := types.AccAddressFromHex(hex.EncodeToString(addressBytes))
	if err != nil {
		return balance, fmt.Errorf("types.AccAddressFromHex: %s", err.Error())
	}
	balance.Available, err = s.node.GetBalance(address.String())
	if err != nil {
		return balance, fmt.Errorf("node.GetBalance: %s", err.Error())
	}
	return balance, nil
}
