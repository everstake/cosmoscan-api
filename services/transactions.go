package services

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/cosmoscan-api/dao/filters"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/everstake/cosmoscan-api/log"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
	"strings"
)

func (s *ServiceFacade) GetAggTransactionsFee(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggTransactionsFee(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggTransactionsFee: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetAggOperationsCount(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAggOperationsCount(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAggOperationsCount: %s", err.Error())
	}
	return items, nil
}

func (s *ServiceFacade) GetAvgOperationsPerBlock(filter filters.Agg) (items []smodels.AggItem, err error) {
	items, err = s.dao.GetAvgOperationsPerBlock(filter)
	if err != nil {
		return nil, fmt.Errorf("dao.GetAvgOperationsPerBlock: %s", err.Error())
	}
	return items, nil
}

type baseMsg struct {
	Type string `json:"@type"`
}

func (s *ServiceFacade) GetTransaction(hash string) (tx smodels.Tx, err error) {
	dTx, err := s.node.GetTransaction(hash)
	if err != nil {
		return tx, fmt.Errorf("node.GetTransaction: %s", err.Error())
	}
	var fee decimal.Decimal
	for _, a := range dTx.Tx.AuthInfo.Fee.Amount {
		if a.Denom == node.MainUnit {
			fee = fee.Add(a.Amount)
		}
	}
	var msgs []smodels.Message
	for _, m := range dTx.Tx.Body.Messages {
		var bm baseMsg
		err = json.Unmarshal(m, &bm)
		if err != nil {
			log.Warn("GetTransaction: parse baseMsg: %s", err.Error())
			continue
		}
		parts := strings.Split(bm.Type, ".")
		t := strings.Trim(parts[len(parts)-1], "Msg")
		msgs = append(msgs, smodels.Message{Type: t, Body: m})
	}
	success := dTx.TxResponse.Code == 0
	fee = node.Precision(fee)
	return smodels.Tx{
		Hash:      dTx.TxResponse.Txhash,
		Type:      dTx.Tx.Type,
		Status:    success,
		Fee:       fee,
		Height:    dTx.TxResponse.Height,
		GasUsed:   dTx.TxResponse.GasUsed,
		GasWanted: dTx.TxResponse.GasWanted,
		Memo:      dTx.Tx.Body.Memo,
		CreatedAt: dmodels.NewTime(dTx.TxResponse.Timestamp),
		Messages:  msgs,
	}, nil
}

func (s *ServiceFacade) GetTransactions(filter filters.Transactions) (resp smodels.PaginatableResponse, err error) {
	dTxs, err := s.dao.GetTransactions(filter)
	if err != nil {
		return resp, fmt.Errorf("dao.GetTransactions: %s", err.Error())
	}
	total, err := s.dao.GetTransactionsCount(filter)
	if err != nil {
		return resp, fmt.Errorf("dao.GetTransactionsCount: %s", err.Error())
	}
	var txs []smodels.TxItem
	for _, tx := range dTxs {
		txs = append(txs, smodels.TxItem{
			Hash:      tx.Hash,
			Status:    tx.Status,
			Fee:       tx.Fee,
			Height:    tx.Height,
			Messages:  tx.Messages,
			CreatedAt: dmodels.NewTime(tx.CreatedAt),
		})
	}
	return smodels.PaginatableResponse{
		Items: txs,
		Total: total,
	}, nil
}
