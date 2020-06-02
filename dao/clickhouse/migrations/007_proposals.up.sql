create table proposals
(
    pro_id           FixedString(42),
    pro_init_deposit Decimal128(18),
    pro_proposer     String,
    pro_content      String,
    pro_created_at   DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMMDD(pro_created_at)
      ORDER BY (pro_id);


