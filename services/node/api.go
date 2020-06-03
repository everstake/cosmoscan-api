package node

import (
	"encoding/json"
	"fmt"
	"github.com/everstake/cosmoscan-api/config"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"net/http"
	"time"
)

const precision = 6

var PrecisionDiv = decimal.New(1, precision)

type (
	API struct {
		cfg    config.Config
		client *http.Client
	}
	CommunityPool struct {
		Height uint64 `json:"height,string"`
		Result [] struct {
			Amount decimal.Decimal `json:"amount"`
		} `json:"result"`
	}
	Validators struct {
		Height uint64      `json:"height,string"`
		Result []Validator `json:"result"`
	}
	Validator struct {
		OperatorAddress string          `json:"operator_address"`
		ConsensusPubkey string          `json:"consensus_pubkey"`
		Jailed          bool            `json:"jailed"`
		Status          int             `json:"status"`
		Tokens          uint64          `json:"tokens,string"`
		DelegatorShares decimal.Decimal `json:"delegator_shares"`
		Description     struct {
			Moniker  string `json:"moniker"`
			Identity string `json:"identity"`
			Website  string `json:"website"`
			Details  string `json:"details"`
		} `json:"description"`
		UnbondingHeight uint64    `json:"unbonding_height,string"`
		UnbondingTime   time.Time `json:"unbonding_time"`
		Commission      struct {
			CommissionRates struct {
				Rate          decimal.Decimal `json:"rate"`
				MaxRate       decimal.Decimal `json:"max_rate"`
				MaxChangeRate decimal.Decimal `json:"max_change_rate"`
			} `json:"commission_rates"`
		} `json:"commission"`
		MaxChangeRate decimal.Decimal `json:"max_change_rate"`
	}
	Inflation struct {
		Height uint64          `json:"height,string"`
		Result decimal.Decimal `json:"result"`
	}
	TotalSupply struct {
		Height uint64 `json:"height,string"`
		Result [] struct {
			Amount decimal.Decimal `json:"amount"`
		} `json:"result"`
	}
	StakingPool struct {
		Height uint64 `json:"height,string"`
		Result struct {
			NotBondedTokens decimal.Decimal `json:"not_bonded_tokens"`
			BondedTokens    decimal.Decimal `json:"bonded_tokens"`
		} `json:"result"`
	}
)

func NewAPI(cfg config.Config) *API {
	return &API{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (api API) request(endpoint string, data interface{}) error {
	url := fmt.Sprintf("%s/%s", api.cfg.Parser.Node, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %s", err.Error())
	}
	resp, err := api.client.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %s", err.Error())
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ioutil.ReadAll: %s", err.Error())
	}
	err = json.Unmarshal(d, data)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %s", err.Error())
	}
	return nil
}

func (api API) GetCommunityPoolAmount() (amount decimal.Decimal, err error) {
	var cp CommunityPool
	err = api.request("distribution/community_pool", &cp)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	if len(cp.Result) == 0 {
		return amount, fmt.Errorf("invalid response")
	}
	amount = cp.Result[0].Amount.Div(PrecisionDiv)
	return amount, nil
}

func (api API) GetValidators() (items []Validator, err error) {
	var validators Validators
	err = api.request("staking/validators", &validators)
	if err != nil {
		return nil, fmt.Errorf("request: %s", err.Error())
	}
	return validators.Result, nil
}

func (api API) GetInflation() (amount decimal.Decimal, err error) {
	var inflation Inflation
	err = api.request("minting/inflation", &inflation)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	return inflation.Result.Mul(decimal.New(100, 0)), nil
}

func (api API) GetTotalSupply() (amount decimal.Decimal, err error) {
	var cp TotalSupply
	err = api.request("supply/total", &cp)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	if len(cp.Result) == 0 {
		return amount, fmt.Errorf("invalid response")
	}
	amount = cp.Result[0].Amount.Div(PrecisionDiv)
	return amount, nil
}

func (api API) GetStakingPool() (sp StakingPool, err error) {
	err = api.request("staking/pool", &sp)
	if err != nil {
		return sp, fmt.Errorf("request: %s", err.Error())
	}
	sp.Result.BondedTokens = sp.Result.BondedTokens.Div(PrecisionDiv)
	sp.Result.NotBondedTokens = sp.Result.NotBondedTokens.Div(PrecisionDiv)
	return sp, nil
}
