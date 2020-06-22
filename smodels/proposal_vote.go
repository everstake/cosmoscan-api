package smodels

import "github.com/everstake/cosmoscan-api/dmodels"

type ProposalVote struct {
	Title       string `json:"title"`
	IsValidator bool   `json:"is_validator"`
	dmodels.ProposalVote
}
