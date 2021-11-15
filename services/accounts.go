package services

import (
	"context"
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/smodels"
	"time"
)

func (s *ServiceFacade) MakeUpdateBalances() {
	tn := time.Now()
	accounts, err := s.dao.GetAccounts(filters.Accounts{})
	if err != nil {
		log.Error("MakeUpdateBalances: dao.GetAccounts: %s", err.Error())
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fetchers := 5
	accountsCh := make(chan dmodels.Account)
	for i := 0; i < fetchers; i++ {
		go func() {
			for {
				select {
				case acc := <-accountsCh:
					for {
						err := s.updateAccount(acc)
						if err != nil {
							log.Warn("MakeSmartUpdateBalances: updateAccount: %s", err.Error())
							time.After(time.Second * 2)
							continue
						}
						break
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
	for _, acc := range accounts {
		accountsCh <- acc
	}
	<-time.After(time.Second * 5)
	log.Info("MakeUpdateBalances finished, duration: %s", time.Now().Sub(tn))
}

func (s *ServiceFacade) updateAccount(account dmodels.Account) error {
	balance, err := s.node.GetBalance(account.Address)
	if err != nil {
		return fmt.Errorf("node.GetBalance: %s", err.Error())
	}
	stake, err := s.node.GetStake(account.Address)
	if err != nil {
		return fmt.Errorf("node.GetStake: %s", err.Error())
	}
	if balance.Equal(account.Balance) && stake.Equal(account.Stake) {
		return nil
	}
	unbonding, err := s.node.GetUnbonding(account.Address)
	if err != nil {
		return fmt.Errorf("node.GetUnbonding: %s", err.Error())
	}
	account.Balance = balance
	account.Stake = stake
	account.Unbonding = unbonding
	err = s.dao.UpdateAccount(account)
	if err != nil {
		return fmt.Errorf("dao.UpdateAccount: %s", err.Error())
	}
	return nil
}

func (s *ServiceFacade) GetAccount(address string) (account smodels.Account, err error) {
	balance, err := s.node.GetBalance(account.Address)
	if err != nil {
		return account, fmt.Errorf("node.GetBalance: %s", err.Error())
	}
	stake, err := s.node.GetStake(account.Address)
	if err != nil {
		return account, fmt.Errorf("node.GetStake: %s", err.Error())
	}
	unbonding, err := s.node.GetUnbonding(account.Address)
	if err != nil {
		return account, fmt.Errorf("node.GetUnbonding: %s", err.Error())
	}
	rewards, err := s.node.GetStakeRewards(account.Address)
	if err != nil {
		return account, fmt.Errorf("node.GetStakeRewards: %s", err.Error())
	}
	return smodels.Account{
		Address:     address,
		Balance:     balance,
		Delegated:   stake,
		Unbonding:   unbonding,
		StakeReward: rewards,
	}, nil
}
