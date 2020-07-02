create table missed_blocks
(
    mib_id          FixedString(40),
    mib_height      UInt64,
    mib_validator   FixedString(40),
    mib_is_proposer UInt8,
    mib_created_at  DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMM(mib_created_at)
      ORDER BY (mib_id);

