CREATE TABLE IF NOT EXISTS account_txs
(
    atx_account FixedString(45),
    atx_tx_hash FixedString(64)
) ENGINE ReplacingMergeTree() ORDER BY (atx_account, atx_tx_hash);

