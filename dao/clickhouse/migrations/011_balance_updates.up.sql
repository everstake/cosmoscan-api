CREATE TABLE IF NOT EXISTS balance_updates
(
    bau_id         FixedString(40),
    bau_address    FixedString(45),
    bau_balance    Decimal(20, 8),
    bau_stake      Decimal(20, 8),
    bau_unbonding  Decimal(20, 8),
    bau_created_at DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMMDD(bau_created_at)
      ORDER BY (bau_id);

