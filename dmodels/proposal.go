package dmodels

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

const ProposalsTable = "proposals"

type Proposal struct {
	ID                uint64          `db:"pro_id" json:"id"`
	TxHash            string          `db:"pro_tx_hash" json:"tx_hash"`
	Type              string          `db:"pro_type" json:"type"`
	Proposer          string          `db:"pro_proposer" json:"proposer"`
	Title             string          `db:"pro_title" json:"title"`
	Description       string          `db:"pro_description" json:"description"`
	Status            string          `db:"pro_status" json:"status"`
	VotesYes          decimal.Decimal `db:"pro_votes_yes" json:"votes_yes"`
	VotesAbstain      decimal.Decimal `db:"pro_votes_abstain" json:"votes_abstain"`
	VotesNo           decimal.Decimal `db:"pro_votes_no" json:"votes_no"`
	VotesNoWithVeto   decimal.Decimal `db:"pro_votes_no_with_veto" json:"votes_no_with_veto"`
	SubmitTime        Time            `db:"pro_submit_time" json:"submit_time"`
	DepositEndTime    Time            `db:"pro_deposit_end_time" json:"deposit_end_time"`
	TotalDeposits     decimal.Decimal `db:"pro_total_deposits" json:"total_deposits"`
	VotingStartTime   Time            `db:"pro_voting_start_time" json:"voting_start_time"`
	VotingEndTime     Time            `db:"pro_voting_end_time" json:"voting_end_time"`
	Voters            uint64          `db:"pro_voters" json:"voters"`
	ParticipationRate decimal.Decimal `db:"pro_participation_rate" json:"participation_rate"`
	Turnout           decimal.Decimal `db:"pro_turnout" json:"turnout"`
	Activity          json.RawMessage `db:"pro_activity" json:"activity"`
}
