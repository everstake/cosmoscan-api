create table delegations
(
    dlg_id         FixedString(40),
    dlg_tx_hash    FixedString(64),
    dlg_delegator  FixedString(46),
    dlg_validator  FixedString(53),
    dlg_amount     Decimal128(18),
    dlg_created_at DateTime
) ENGINE ReplacingMergeTree()
PARTITION BY toYYYYMMDD(dlg_created_at)
ORDER BY (dlg_id);
