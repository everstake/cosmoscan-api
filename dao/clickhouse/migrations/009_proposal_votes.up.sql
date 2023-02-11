CREATE TABLE IF NOT EXISTS proposal_votes
(
    prv_id          FixedString(40),
    prv_proposal_id UInt64,
    prv_tx_hash     FixedString(64),
    prv_voter       String,
    prv_option      String,
    prv_created_at  DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMMDD(prv_created_at)
      ORDER BY (prv_id);
