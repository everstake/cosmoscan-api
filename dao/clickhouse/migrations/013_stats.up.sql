CREATE TABLE IF NOT EXISTS stats
(
    stt_id         FixedString(42),
    stt_title      String,
    stt_value      String,
    stt_created_at DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMM(stt_created_at)
      ORDER BY (stt_id);
