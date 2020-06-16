package services

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services/parser/hub3"
	"github.com/shopspring/decimal"
	"time"
)

func (s *ServiceFacade) MakeUpdateBalances() {
	tn := time.Now()
	updates, err := s.dao.GetBalanceUpdate(filters.BalanceUpdates{Limit: 1})
	if err != nil {
		log.Error("MakeUpdateBalances: dao.GetBalanceUpdate: %s", err.Error())
		return
	}
	if len(updates) == 0 {
		err = s.makeFirstUpdateBalance()
		if err != nil {
			log.Error("MakeUpdateBalances: makeFirstUpdateBalance: %s", err.Error())
			return
		}
	} else {
		// regular update
		activeAccounts, err := s.dao.GetActiveAccounts(filters.ActiveAccounts{
			From: time.Now().Add(time.Hour * 24 * 2),
		})
		if err != nil {
			log.Error("MakeUpdateBalances: dao.GetActiveAccounts: %s", err.Error())
			return
		}
		s.makeUpdateBalance(activeAccounts)
	}
	log.Info("MakeUpdateBalances finished, duration: %s", time.Now().Sub(tn))
}

func (s *ServiceFacade) makeFirstUpdateBalance() error {
	// check if data has been already parsed && wait
	for {
		parser, err := s.dao.GetParser(hub3.ParserTitle)
		if err != nil {
			log.Warn("makeFirstUpdateBalance: dao.GetParser: %s", err.Error())
			<-time.After(time.Second * 5)
			continue
		}
		node := hub3.NewAPI(s.cfg.Parser.Node)
		block, err := node.GetLatestBlock()
		if err != nil {
			log.Warn("makeFirstUpdateBalance: node.GetLatestBlock: %s", err.Error())
			<-time.After(time.Second * 5)
			continue
		}
		diff := block.Block.Header.Height - parser.Height
		if diff < 0 {
			diff = - diff
		}
		if diff > 100 {
			<-time.After(time.Minute)
		} else {
			break
		}
	}
	activeAccounts, err := s.dao.GetActiveAccounts(filters.ActiveAccounts{})
	if err != nil {
		return fmt.Errorf("dao.GetActiveAccounts: %s", err.Error())
	}
	fmt.Println("activeAccounts", len(activeAccounts))
	allAccounts, err := s.dao.GetAccounts(filters.Accounts{})
	if err != nil {
		return fmt.Errorf("dao.GetAccounts: %s", err.Error())
	}
	fmt.Println("allAccounts", len(allAccounts))

	allUniqAccounts := make(map[string]struct{})
	for _, acc := range activeAccounts {
		allUniqAccounts[acc] = struct{}{}
	}
	for _, acc := range allAccounts {
		allUniqAccounts[acc.Address] = struct{}{}
	}
	var addresses []string
	for address := range allUniqAccounts {
		addresses = append(addresses, address)
	}
	s.makeUpdateBalance(addresses)
	return nil
}

func (s *ServiceFacade) makeUpdateBalance(accounts []string) {
	var err error
	for _, acc := range accounts {
		var balance, stake, unbonding decimal.Decimal
		for {
			balance, err = s.node.GetBalance(acc)
			if err == nil {
				break
			}
			log.Warn("makeUpdateBalance: node.GetBalance: %s", err.Error())
			<-time.After(time.Second * 5)
		}
		for {
			stake, err = s.node.GetStake(acc)
			if err == nil {
				break
			}
			log.Warn("makeUpdateBalance: node.GetStake: %s", err.Error())
			<-time.After(time.Second * 5)
		}
		for {
			unbonding, err = s.node.GetUnbonding(acc)
			if err == nil {
				break
			}
			log.Warn("makeUpdateBalance: node.GetUnbonding: %s", err.Error())
			<-time.After(time.Second * 5)
		}
		balance = balance.Truncate(8)
		stake = stake.Truncate(8)
		unbonding = unbonding.Truncate(8)
		account, err := s.dao.GetAccount(acc)
		if err != nil {
			log.Error("makeUpdateBalance: dao.GetAccount: %s", err.Error())
			continue
		}
		if account.Stake.Equal(stake) || account.Balance.Equal(balance) {
			continue
		}
		account.Balance = balance
		account.Stake = stake
		account.Unbonding = unbonding
		err = s.dao.UpdateAccount(account)
		if err != nil {
			log.Error("makeUpdateBalance: dao.UpdateAccount: %s", err.Error())
			continue
		}
		tn := time.Now()
		hash := sha1.Sum([]byte(fmt.Sprintf("%s.%s", account.Address, tn.String())))
		id := hex.EncodeToString(hash[:])
		err = s.dao.CreateBalanceUpdates([]dmodels.BalanceUpdate{
			{
				ID:        id,
				Address:   account.Address,
				Stake:     account.Stake,
				Balance:   account.Balance,
				Unbonding: account.Unbonding,
				CreatedAt: tn,
			},
		})
		if err != nil {
			log.Error("makeUpdateBalance: dao.CreateBalanceUpdates: %s", err.Error())
			continue
		}
	}
}
