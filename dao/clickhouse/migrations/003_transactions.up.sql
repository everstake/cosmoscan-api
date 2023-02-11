CREATE TABLE IF NOT EXISTS transactions
(
    trn_hash       FixedString(64),
    trn_block_id   UInt64,
    trn_status     UInt8,
    trn_height     UInt64,
    trn_messages   UInt32,
    trn_fee        Decimal128(18),
    trn_gas_used   UInt64,
    trn_gas_wanted UInt64,
    trn_created_at DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMMDD(trn_created_at)
      ORDER BY (trn_hash);
