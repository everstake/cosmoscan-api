create table delegator_rewards
(
    der_id         FixedString(40),
    der_tx_hash    FixedString(64),
    der_delegator  FixedString(50),
    der_validator  FixedString(57),
    der_amount     Decimal128(18),
    der_created_at DateTime
) ENGINE ReplacingMergeTree()
PARTITION BY toYYYYMMDD(der_created_at)
ORDER BY (der_id);