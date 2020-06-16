package dmodels

import "time"

const JailersTable = "jailers"

type Jailer struct {
	ID        string    `db:"jlr_id"`
	Address   string    `db:"jlr_address"`
	CreatedAt time.Time `db:"jlr_created_at"`
}
