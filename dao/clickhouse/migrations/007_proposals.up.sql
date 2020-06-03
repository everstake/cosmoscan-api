create table proposals
(
    pro_id           UInt64,
    pro_title        String,
    pro_description  String,
    pro_recipient    String,
    pro_amount       Decimal128(18),
    pro_init_deposit Decimal128(18),
    pro_proposer     String,
    pro_created_at   DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMMDD(pro_created_at)
      ORDER BY (pro_id);
