package services

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/cosmoscan-api/services/node"
	"github.com/everstake/cosmoscan-api/smodels"
	"github.com/shopspring/decimal"
)

type (
	baseMsg struct {
		Type string `json:"@type"`
	}
	msgAmount struct {
		Denom  string          `json:"denom"`
		Amount decimal.Decimal `json:"amount"`
	}
	sendMsg struct {
		FromAddress string      `json:"from_address"`
		ToAddress   string      `json:"to_address"`
		Amount      []msgAmount `json:"amount"`
	}
)

func parseMsg(msg json.RawMessage) (model interface{}, err error) {
	var bMsg baseMsg
	err = json.Unmarshal(msg, &bMsg)
	if err != nil {
		return model, fmt.Errorf("json.Unmarshal(baseMsg): %s", err.Error())
	}
	switch bMsg.Type {
	case "/cosmos.bank.v1beta1.MsgSend":
		m := sendMsg{}
		err = json.Unmarshal(msg, &m)
		if err != nil {
			return model, fmt.Errorf("json.Unmarshal(%s): %s", bMsg.Type, err.Error())
		}
		model = smodels.SendMsg{
			From:   m.FromAddress,
			To:     m.ToAddress,
			Amount: calculateMainSum(m.Amount),
		}
	}
	return nil, nil
}

func calculateMainSum(amount []msgAmount) (sum decimal.Decimal) {
	for _, a := range amount {
		if a.Denom == node.MainUnit {
			sum = sum.Add(a.Amount)
		}
	}
	return sum.Div(node.PrecisionDiv)
}
