package dmodels

const AccountTxsTable = "account_txs"

type AccountTx struct {
	Account string `db:"atx_account"`
	TxHash  string `db:"atx_tx_hash"`
}
