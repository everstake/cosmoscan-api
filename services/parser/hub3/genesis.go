package hub3

import "github.com/shopspring/decimal"

const genesisJson = "https://raw.githubusercontent.com/cosmos/launch/master/genesis.json"

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
				ValidatorAddress decimal.Decimal `json:"validator_address"`
			} `json:"delegations"`
		} `json:"staking"`
	} `json:"app_state"`
}
