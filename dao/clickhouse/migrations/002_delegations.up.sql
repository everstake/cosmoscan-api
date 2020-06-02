create table delegations
(
    dlg_id         FixedString(42),
    dlg_tx_hash    FixedString(64),
    dlg_delegator  FixedString(45),
    dlg_validator  FixedString(52),
    dlg_amount     Decimal128(18),
    dlg_created_at DateTime
) ENGINE ReplacingMergeTree()
PARTITION BY toYYYYMMDD(dlg_created_at)
ORDER BY (dlg_id);
