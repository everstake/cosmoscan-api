package dmodels

import "time"

const MissedBlocks = "missed_blocks"

type MissedBlock struct {
	ID         string    `db:"mib_id"`
	Height     uint64    `db:"mib_height"`
	Validator  string    `db:"mib_validator"`
	IsProposer bool      `db:"mib_is_proposer"`
	CreatedAt  time.Time `db:"mib_created_at"`
}
