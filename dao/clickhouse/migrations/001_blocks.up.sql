CREATE TABLE IF NOT EXISTS blocks
(
    blk_id         UInt64,
    blk_hash       FixedString(64),
    blk_proposer   FixedString(40),
    blk_created_at DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMMDD(blk_created_at)
      ORDER BY (blk_id);
