package services

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
	"math"
	"sort"
	"time"
)

func (s *ServiceFacade) GetNetworkStates(filter filters.Stats) (map[string][]decimal.Decimal, error) {
	filter.Titles = []string{
		dmodels.StatsTotalStakingBalance,
		dmodels.StatsNumberDelegators,
		dmodels.StatsNumberMultiDelegators,
		dmodels.StatsTransfersVolume,
		dmodels.StatsFeeVolume,
		dmodels.StatsHighestFee,
		dmodels.StatsUndelegationVolume,
		dmodels.StatsBlockDelay,
		dmodels.StatsNetworkSize,
		dmodels.StatsTotalAccounts,
		dmodels.StatsTotalWhaleAccounts,
		dmodels.StatsTotalSmallAccounts,
		dmodels.StatsTotalJailers,
	}
	return s.getStates(filter)
}

func (s *ServiceFacade) getStates(filter filters.Stats) (map[string][]decimal.Decimal, error) {
	if filter.To.IsZero() {
		filter.To = dmodels.NewTime(time.Now())
	}
	filter.From = dmodels.NewTime(filter.To.Add(-time.Hour * 24 * 7))
	stats, err := s.dao.GetStats(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetStats: %s", err.Error())
	}
	mp := make(map[string][]decimal.Decimal)
	for _, stat := range stats {
		mp[stat.Title] = append(mp[stat.Title], stat.Value)
	}
	return mp, nil
}

func (s *ServiceFacade) MakeStats() {
	now := time.Now()
	y, m, d := now.Add(-time.Hour * 24).Date()
	startOfYesterday := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	startOfToday := startOfYesterday.Add(time.Hour * 24)
	stats := []struct {
		title string
		fetch func() (decimal.Decimal, error)
	}{
		{
			title: dmodels.StatsTotalStakingBalance,
			fetch: func() (value decimal.Decimal, err error) {
				stakingPool, err := s.node.GetStakingPool()
				if err != nil {
					return value, fmt.Errorf("node.GetStakingPool: %s", err.Error())
				}
				return stakingPool.Result.BondedTokens, nil
			},
		},
		{
			title: dmodels.StatsNumberDelegators,
			fetch: func() (value decimal.Decimal, err error) {
				total, err := s.dao.GetDelegatorsTotal(filters.Delegators{
					TimeRange: filters.TimeRange{
						From: dmodels.NewTime(startOfYesterday),
						To:   dmodels.NewTime(startOfToday),
					},
				})
				if err != nil {
					return value, fmt.Errorf("dao.GetDelegatorsTotal: %s", err.Error())
				}
				return decimal.NewFromInt(int64(total)), nil
			},
		},
		{
			title: dmodels.StatsNumberMultiDelegators,
			fetch: func() (value decimal.Decimal, err error) {
				total, err := s.dao.GetMultiDelegatorsTotal(filters.TimeRange{})
				if err != nil {
					return value, fmt.Errorf("dao.GetMultiDelegatorsTotal: %s", err.Error())
				}
				return decimal.NewFromInt(int64(total)), nil
			},
		},
		{
			title: dmodels.StatsTransfersVolume,
			fetch: func() (value decimal.Decimal, err error) {
				volume, err := s.dao.GetTransferVolume(filters.TimeRange{})
				if err != nil {
					return value, fmt.Errorf("dao.GetTransferVolume: %s", err.Error())
				}
				return volume, nil
			},
		},
		{
			title: dmodels.StatsFeeVolume,
			fetch: func() (value decimal.Decimal, err error) {
				volume, err := s.dao.GetTransactionsFeeVolume(filters.TimeRange{
					From: dmodels.NewTime(startOfYesterday),
					To:   dmodels.NewTime(startOfToday),
				})
				if err != nil {
					return value, fmt.Errorf("dao.GetTransactionsFeeVolume: %s", err.Error())
				}
				return volume, nil
			},
		},
		{
			title: dmodels.StatsHighestFee,
			fetch: func() (value decimal.Decimal, err error) {
				volume, err := s.dao.GetTransactionsHighestFee(filters.TimeRange{
					From: dmodels.NewTime(startOfYesterday),
					To:   dmodels.NewTime(startOfToday),
				})
				if err != nil {
					return value, fmt.Errorf("dao.GetTransactionsHighestFee: %s", err.Error())
				}
				return volume, nil
			},
		},
		{
			title: dmodels.StatsUndelegationVolume,
			fetch: func() (value decimal.Decimal, err error) {
				volume, err := s.dao.GetUndelegationsVolume(filters.TimeRange{
					From: dmodels.NewTime(startOfYesterday),
					To:   dmodels.NewTime(startOfToday),
				})
				if err != nil {
					return value, fmt.Errorf("dao.GetUndelegationsVolume: %s", err.Error())
				}
				return volume, nil
			},
		},
		{
			title: dmodels.StatsBlockDelay,
			fetch: func() (value decimal.Decimal, err error) {
				delay, err := s.dao.GetAvgBlocksDelay(filters.TimeRange{
					From: dmodels.NewTime(startOfYesterday),
					To:   dmodels.NewTime(startOfToday),
				})
				if err != nil {
					return value, fmt.Errorf("dao.GetAvgBlocksDelay: %s", err.Error())
				}
				if math.IsNaN(delay) {
					return decimal.Zero, nil
				}
				return decimal.NewFromInt(int64(delay)), nil
			},
		},
		{
			title: dmodels.StatsNetworkSize,
			fetch: func() (value decimal.Decimal, err error) {
				size, err := s.GetSizeOfNode()
				if err != nil {
					return value, fmt.Errorf("GetSizeOfNode: %s", err.Error())
				}
				return decimal.NewFromFloat(size), nil
			},
		},
		{
			title: dmodels.StatsTotalAccounts,
			fetch: func() (value decimal.Decimal, err error) {
				total, err := s.dao.GetAccountsTotal(filters.Accounts{})
				if err != nil {
					return value, fmt.Errorf("dao.GetAccountsTotal: %s", err.Error())
				}
				return decimal.NewFromInt(int64(total)), nil
			},
		},
		{
			title: dmodels.StatsTotalWhaleAccounts,
			fetch: func() (value decimal.Decimal, err error) {
				minAmount := decimal.NewFromFloat(300000)
				total, err := s.dao.GetAccountsTotal(filters.Accounts{GtTotalAmount: minAmount})
				if err != nil {
					return value, fmt.Errorf("dao.GetAccountsTotal: %s", err.Error())
				}
				return decimal.NewFromInt(int64(total)), nil
			},
		},
		{
			title: dmodels.StatsTotalSmallAccounts,
			fetch: func() (value decimal.Decimal, err error) {
				maxAmount := decimal.NewFromFloat(1)
				total, err := s.dao.GetAccountsTotal(filters.Accounts{LtTotalAmount: maxAmount})
				if err != nil {
					return value, fmt.Errorf("dao.GetAccountsTotal: %s", err.Error())
				}
				return decimal.NewFromInt(int64(total)), nil
			},
		},
		{
			title: dmodels.StatsTotalJailers,
			fetch: func() (value decimal.Decimal, err error) {
				total, err := s.dao.GetJailersTotal()
				if err != nil {
					return value, fmt.Errorf("dao.GetJailersTotal: %s", err.Error())
				}
				return decimal.NewFromInt(int64(total)), nil
			},
		},
		{
			title: dmodels.StatsValidatorsWith33Power,
			fetch: func() (value decimal.Decimal, err error) {
				sp, err := s.node.GetStakingPool()
				if err != nil {
					return value, fmt.Errorf("node.GetStakingPool: %s", err.Error())
				}
				mp, err := s.GetValidatorMap()
				if err != nil {
					return value, fmt.Errorf("s.GetValidatorMap: %s", err.Error())
				}
				var amounts []decimal.Decimal
				for _, validator := range mp {
					amounts = append(amounts, validator.DelegatorShares.Div(node.PrecisionDiv))
				}
				sort.Slice(amounts, func(i, j int) bool {
					return amounts[i].GreaterThan(amounts[j])
				})
				if sp.Result.BondedTokens.IsZero() {
					return value, fmt.Errorf("total stake is zero")
				}
				stake := sp.Result.BondedTokens
				sum := decimal.Zero
				limit := decimal.NewFromFloat(33.4)
				for _, amount := range amounts {
					sum = sum.Add(amount)
					value = value.Add(decimal.NewFromInt(1))
					power := sum.Div(stake).Mul(decimal.NewFromInt(100))
					if power.GreaterThan(limit) {
						return value, nil
					}
				}
				return value, nil
			},
		},
	}

	var models []dmodels.Stat
	for _, stat := range stats {
		value, err := stat.fetch()
		if err != nil {
			log.Error("MakeStats (%s): %s", stat.title, err.Error())
			continue
		}
		hash := sha1.Sum([]byte(fmt.Sprintf("%s.%s", stat.title, startOfToday.String())))
		id := hex.EncodeToString(hash[:])
		models = append(models, dmodels.Stat{
			ID:        id,
			Title:     stat.title,
			Value:     value,
			CreatedAt: startOfToday,
		})
	}
	err := s.dao.CreateStats(models)
	if err != nil {
		log.Error("MakeStats: dao.CreateStats: %s", err.Error())
	}
}

func (s *ServiceFacade) GetAggValidators33Power(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggValidators33Power(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggValidators33Power: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetAggWhaleAccounts(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggWhaleAccounts(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggWhaleAccounts: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetAggBondedRatio(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggHistoricalStatesByField(filter, "his_staked_ratio")
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggHistoricalStatesByField: %s", err.Error())
	}
	return items, nil
}
