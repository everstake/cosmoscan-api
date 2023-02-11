CREATE TABLE IF NOT EXISTS validator_rewards
(
    var_id         FixedString(40),
    var_tx_hash    FixedString(64),
    var_address    FixedString(52),
    var_amount     Decimal128(18),
    var_created_at DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMMDD(var_created_at)
      ORDER BY (var_id);
