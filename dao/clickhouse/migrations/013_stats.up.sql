create table stats
(
    stt_id         FixedString(42),
    stt_title      String,
    stt_created_at DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMM(stt_created_at)
      ORDER BY (stt_id);
