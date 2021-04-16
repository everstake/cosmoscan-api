package hub3

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/cosmoscan-api/dmodels"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"time"
)

const genesisJson = "https://raw.githubusercontent.com/cosmos/launch/master/genesis.json"
const saveGenesisBatch = 100

type Genesis struct {
	AppState struct {
		Accounts []struct {
			Address string   `json:"address"`
			Coins   []Amount `json:"coins"`
		} `json:"accounts"`
		Distribution struct {
			DelegatorStartingInfos []struct {
				StartingInfo struct {
					DelegatorAddress string `json:"delegator_address"`
					StartingInfo     struct {
						Stake decimal.Decimal `json:"stake"`
					} `json:"starting_info"`
					ValidatorAddress string `json:"validator_address"`
				} `json:"starting_info"`
			} `json:"delegator_starting_infos"`
		} `json:"distribution"`
		Staking struct {
			Delegations []struct {
				DelegatorAddress string          `json:"delegator_address"`
				Shares           decimal.Decimal `json:"shares"`
				ValidatorAddress string          `json:"validator_address"`
			} `json:"delegations"`
			Redelegations []struct {
				DelegatorAddress string `json:"delegator_address"`
				Entries          []struct {
					SharesDst decimal.Decimal `json:"shares_dst"`
				} `json:"entries"`
				ValidatorDstAddress string `json:"validator_dst_address"`
				ValidatorSrcAddress string `json:"validator_src_address"`
			} `json:"redelegations"`
		} `json:"staking"`
	} `json:"app_state"`
	GenesisTime time.Time `json:"genesis_time"`
	Validators  []struct {
		Address string          `json:"address"`
		Name    string          `json:"name"`
		Power   decimal.Decimal `json:"power"`
	} `json:"validators"`
}

func GetGenesisState() (state Genesis, err error) {
	resp, err := http.Get(genesisJson)
	if err != nil {
		return state, fmt.Errorf("http.Get: %s", err.Error())
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return state, fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	err = json.Unmarshal(data, &state)
	if err != nil {
		return state, fmt.Errorf("json.Unmarshal: %s", err.Error())
	}

	return state, nil
}

func ShowGenesisStructure() {
	resp, _ := http.Get(genesisJson)
	data, _ := ioutil.ReadAll(resp.Body)
	var value interface{}
	_ = json.Unmarshal(data, &value)
	printStruct(value, 0)
}

func printStruct(field interface{}, i int) {
	mp, ok := field.(map[string]interface{})
	if ok {
		if len(mp) > 50 {
			return
		}
		for title, f := range mp {
			var str string
			for k := 0; k < i; k++ {
				str = str + " "
			}
			fmt.Println(str + title)
			printStruct(f, i+1)
		}
	}
}

func (p *Parser) parseGenesisState() error {
	state, err := GetGenesisState()
	if err != nil {
		return fmt.Errorf("getGenesisState: %s", err.Error())
	}
	t, err := time.Parse("2006-01-02", "2019-12-11")
	if err != nil {
		return fmt.Errorf("time.Parse: %s", err.Error())
	}
	var (
		delegations []dmodels.Delegation
		accounts    []dmodels.Account
	)
	for i, delegation := range state.AppState.Staking.Delegations {
		delegations = append(delegations, dmodels.Delegation{
			ID:        makeHash(fmt.Sprintf("delegations.%d", i)),
			TxHash:    "genesis",
			Delegator: delegation.DelegatorAddress,
			Validator: delegation.ValidatorAddress,
			Amount:    delegation.Shares.Div(precisionDiv),
			CreatedAt: t,
		})
	}
	for i, delegation := range state.AppState.Staking.Redelegations {
		amount := decimal.Zero
		for _, entry := range delegation.Entries {
			amount = amount.Add(entry.SharesDst)
		}
		// ignore undelegation
		delegations = append(delegations, dmodels.Delegation{
			ID:        makeHash(fmt.Sprintf("redelegations.%d", i)),
			TxHash:    "genesis",
			Delegator: delegation.DelegatorAddress,
			Validator: delegation.ValidatorDstAddress,
			Amount:    amount.Div(precisionDiv),
			CreatedAt: t,
		})
	}
	accountDelegation := make(map[string]decimal.Decimal)
	for _, delegation := range delegations {
		accountDelegation[delegation.Delegator] = accountDelegation[delegation.Delegator].Add(delegation.Amount)
	}
	for _, account := range state.AppState.Accounts {
		amount, _ := calculateAtomAmount(account.Coins)
		accounts = append(accounts, dmodels.Account{
			Address:   account.Address,
			Balance:   amount,
			Stake:     accountDelegation[account.Address],
			CreatedAt: t,
		})
	}

	for i := 0; i < len(accounts); i += saveGenesisBatch {
		endOfPart := i + saveGenesisBatch
		if i+saveGenesisBatch > len(accounts) {
			endOfPart = len(accounts)
		}
		err := p.dao.CreateAccounts(accounts[i:endOfPart])
		if err != nil {
			return fmt.Errorf("dao.CreateAccounts: %s", err.Error())
		}
	}

	for i := 0; i < len(delegations); i += saveGenesisBatch {
		endOfPart := i + saveGenesisBatch
		if i+saveGenesisBatch > len(delegations) {
			endOfPart = len(delegations)
		}
		err := p.dao.CreateDelegations(delegations[i:endOfPart])
		if err != nil {
			return fmt.Errorf("dao.CreateDelegations: %s", err.Error())
		}
	}

	return nil
}
