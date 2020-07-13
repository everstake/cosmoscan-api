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
	AmountResult struct {
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
	StakeResult struct {
		Height uint64 `json:"height,string"`
		Result []struct {
			DelegatorAddress string          `json:"delegator_address"`
			ValidatorAddress string          `json:"validator_address"`
			Shares           decimal.Decimal `json:"shares"`
		} `json:"result"`
	}
	UnbondingResult struct {
		Height uint64 `json:"height,string"`
		Result []struct {
			DelegatorAddress string `json:"delegator_address"`
			ValidatorAddress string `json:"validator_address"`
			Entries          [] struct {
				Balance decimal.Decimal `json:"balance"`
			} `json:"entries"`
		} `json:"result"`
	}
	ProposalsResult struct {
		Height uint64 `json:"height,string"`
		Result []struct {
			Content struct {
				Type  string `json:"type"`
				Value struct {
					Title       string `json:"title"`
					Description string `json:"description"`
				} `json:"value"`
			} `json:"content"`
			ID               uint64 `json:"id,string"`
			ProposalStatus   string `json:"proposal_status"`
			FinalTallyResult struct {
				Yes        int64 `json:"yes,string"`
				Abstain    int64 `json:"abstain,string"`
				No         int64 `json:"no,string"`
				NoWithVeto int64 `json:"no_with_veto,string"`
			} `json:"final_tally_result"`
			SubmitTime     time.Time `json:"submit_time"`
			DepositEndTime time.Time `json:"deposit_end_time"`
			TotalDeposit   [] struct {
				Amount decimal.Decimal `json:"amount"`
			} `json:"total_deposit"`
			VotingStartTime time.Time `json:"voting_start_time"`
			VotingEndTime   time.Time `json:"voting_end_time"`
		} `json:"result"`
	}
	ProposalProposer struct {
		Height uint64 `json:"height,string"`
		Result struct {
			ProposalID uint64 `json:"proposal_id,string"`
			Proposer   string `json:"proposer"`
		} `json:"result"`
	}
	DelegatorValidatorStakeResult struct {
		Height uint64 `json:"height,string"`
		Result struct {
			DelegatorAddress string          `json:"delegator_address"`
			ValidatorAddress string          `json:"validator_address"`
			Shares           decimal.Decimal `json:"shares"`
			Balance          decimal.Decimal `json:"balance"`
		} `json:"result"`
	}
	ProposalVotersResult struct {
		Result []struct {
			ProposalID uint64 `json:"proposal_id,string"`
			Voter      string `json:"voter"`
			Option     string `json:"option"`
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
	var cp AmountResult
	err = api.request("supply/total", &cp)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	if len(cp.Result) == 0 {
		return amount, fmt.Errorf("invalid response")
	}
	return cp.Amount(), nil
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

func (api API) GetBalance(address string) (amount decimal.Decimal, err error) {
	var result AmountResult
	err = api.request(fmt.Sprintf("/bank/balances/%s", address), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	return result.Amount(), nil
}

func (api API) GetStake(address string) (amount decimal.Decimal, err error) {
	var result StakeResult
	err = api.request(fmt.Sprintf("/staking/delegators/%s/delegations", address), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	shares := decimal.Zero
	for _, r := range result.Result {
		shares = shares.Add(r.Shares)
	}
	return shares.Div(PrecisionDiv), nil
}

func (api API) GetUnbonding(address string) (amount decimal.Decimal, err error) {
	var result UnbondingResult
	err = api.request(fmt.Sprintf("/staking/delegators/%s/unbonding_delegations", address), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	for _, r := range result.Result {
		for _, entry := range r.Entries {
			amount = amount.Add(entry.Balance)
		}
	}
	amount = amount.Div(PrecisionDiv)
	return amount, nil
}

func (api API) GetProposals() (proposals ProposalsResult, err error) {
	err = api.request("/gov/proposals", &proposals)
	if err != nil {
		return proposals, fmt.Errorf("request: %s", err.Error())
	}
	return proposals, nil
}

func (api API) GetProposalProposer(id uint64) (address string, err error) {
	var result ProposalProposer
	err = api.request(fmt.Sprintf("/gov/proposals/%d/proposer", id), &result)
	if err != nil {
		return address, fmt.Errorf("request: %s", err.Error())
	}
	return result.Result.Proposer, nil
}

func (api API) GetDelegatorValidatorStake(delegator string, validator string) (amount decimal.Decimal, err error) {
	var result DelegatorValidatorStakeResult
	err = api.request(fmt.Sprintf("/staking/delegators/%s/delegations/%s", delegator, validator), &result)
	if err != nil {
		return amount, fmt.Errorf("request: %s", err.Error())
	}
	return result.Result.Shares.Div(PrecisionDiv), nil
}

func (api API) GetProposalVoters(id uint64) (result ProposalVotersResult, err error) {
	err = api.request(fmt.Sprintf("/gov/proposals/%d/votes", id), &result)
	if err != nil {
		return result, fmt.Errorf("request: %s", err.Error())
	}
	return result, nil
}

func (r *AmountResult) Amount() decimal.Decimal {
	s := decimal.Zero
	for _, a := range r.Result {
		s = s.Add(a.Amount)
	}
	return s.Div(PrecisionDiv)
}
