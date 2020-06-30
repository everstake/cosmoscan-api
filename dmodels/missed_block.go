package dmodels

import "time"

const MissedBlocks = "missed_blocks"

type MissedBlock struct {
	ID        string    `db:"mib_id"`
	Height    uint64    `db:"mib_height"`
	Validator string    `db:"mib_validator"`
	CreatedAt time.Time `db:"mib_created_at"`
}
