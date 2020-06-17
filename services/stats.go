package services

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/shopspring/decimal"
	"math"
	"time"
)

const (
	StatsTotalStakingBalance   = "total_staking_balance"
	StatsNumberDelegators      = "number_delegators"
	StatsNumberMultiDelegators = "number_multi_delegators"
	StatsTransfersVolume       = "transfer_volume"
	StatsFeeVolume             = "fee_volume"
	StatsHighestFee            = "highest_fee"
	StatsUndelegationVolume    = "undelegation_volume"
	StatsBlockDelay            = "block_delay"
	StatsNetworkSize           = "network_size"
	StatsTotalAccounts         = "total_accounts"
	StatsTotalWhaleAccounts    = "total_whale_accounts"
	StatsTotalSmallAccounts    = "total_small_accounts"
	StatsTotalJailers          = "total_jailers"
)

func (s *ServiceFacade) GetNetworkStates(filter filters.Stats) (map[string][]decimal.Decimal, error) {
	filter.Titles = []string{
		StatsTotalStakingBalance,
		StatsNumberDelegators,
		StatsNumberMultiDelegators,
		StatsTransfersVolume,
		StatsFeeVolume,
		StatsHighestFee,
		StatsUndelegationVolume,
		StatsBlockDelay,
		StatsNetworkSize,
		StatsTotalAccounts,
		StatsTotalWhaleAccounts,
		StatsTotalSmallAccounts,
		StatsTotalJailers,
	}
	return s.getStates(filter)
}

func (s *ServiceFacade) getStates(filter filters.Stats) (map[string][]decimal.Decimal, error) {
	if filter.To.IsZero() {
		filter.To = dmodels.NewTime(time.Now())
	}
	filter.From = dmodels.NewTime(filter.To.Add(time.Hour * 24 * 7))
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
			title: StatsTotalStakingBalance,
			fetch: func() (value decimal.Decimal, err error) {
				stakingPool, err := s.node.GetStakingPool()
				if err != nil {
					return value, fmt.Errorf("node.GetStakingPool: %s", err.Error())
				}
				return stakingPool.Result.BondedTokens, nil
			},
		},
		{
			title: StatsNumberDelegators,
			fetch: func() (value decimal.Decimal, err error) {
				total, err := s.dao.GetDelegatorsTotal(filters.TimeRange{
					From: dmodels.NewTime(startOfYesterday),
					To:   dmodels.NewTime(startOfToday),
				})
				if err != nil {
					return value, fmt.Errorf("dao.GetDelegatorsTotal: %s", err.Error())
				}
				return decimal.NewFromInt(int64(total)), nil
			},
		},
		{
			title: StatsNumberMultiDelegators,
			fetch: func() (value decimal.Decimal, err error) {
				total, err := s.dao.GetMultiDelegatorsTotal(filters.TimeRange{
					From: dmodels.NewTime(startOfYesterday),
					To:   dmodels.NewTime(startOfToday),
				})
				if err != nil {
					return value, fmt.Errorf("dao.GetMultiDelegatorsTotal: %s", err.Error())
				}
				return decimal.NewFromInt(int64(total)), nil
			},
		},
		{
			title: StatsTransfersVolume,
			fetch: func() (value decimal.Decimal, err error) {
				volume, err := s.dao.GetTransferVolume(filters.TimeRange{
					From: dmodels.NewTime(startOfYesterday),
					To:   dmodels.NewTime(startOfToday),
				})
				if err != nil {
					return value, fmt.Errorf("dao.GetTransferVolume: %s", err.Error())
				}
				return volume, nil
			},
		},
		{
			title: StatsFeeVolume,
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
			title: StatsHighestFee,
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
			title: StatsUndelegationVolume,
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
			title: StatsBlockDelay,
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
			title: StatsNetworkSize,
			fetch: func() (value decimal.Decimal, err error) {
				size, err := s.GetSizeOfNode()
				if err != nil {
					return value, fmt.Errorf("GetSizeOfNode: %s", err.Error())
				}
				return decimal.NewFromFloat(size), nil
			},
		},
		{
			title: StatsTotalAccounts,
			fetch: func() (value decimal.Decimal, err error) {
				total, err := s.dao.GetAccountsTotal(filters.Accounts{})
				if err != nil {
					return value, fmt.Errorf("dao.GetAccountsTotal: %s", err.Error())
				}
				return decimal.NewFromInt(int64(total)), nil
			},
		},
		{
			title: StatsTotalWhaleAccounts,
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
			title: StatsTotalSmallAccounts,
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
			title: StatsTotalJailers,
			fetch: func() (value decimal.Decimal, err error) {
				total, err := s.dao.GetJailersTotal()
				if err != nil {
					return value, fmt.Errorf("dao.GetJailersTotal: %s", err.Error())
				}
				return decimal.NewFromInt(int64(total)), nil
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
