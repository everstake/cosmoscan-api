CREATE TABLE IF NOT EXISTS proposal_deposits
(
    prd_id          FixedString(40),
    prd_proposal_id UInt64,
    prd_depositor   String,
    prd_amount      Decimal128(18),
    prd_created_at  DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMMDD(prd_created_at)
      ORDER BY (prd_id);

