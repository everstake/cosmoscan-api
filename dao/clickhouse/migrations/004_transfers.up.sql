create table transfers
(
    trf_id         FixedString(40),
    trf_tx_hash    FixedString(64),
    trf_from       FixedString(45),
    trf_to         FixedString(45),
    trf_amount     Decimal128(18),
    trf_created_at DateTime,
    trf_currency   String
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMMDD(trf_created_at)
      ORDER BY (trf_id);
