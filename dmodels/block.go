package dmodels

import "time"

const BlocksTable = "blocks"

type Block struct {
	ID        uint64    `db:"blk_id"`
	Hash      string    `db:"blk_hash"`
	Proposer  string    `db:"blk_proposer"`
	CreatedAt time.Time `db:"blk_created_at"`
}
