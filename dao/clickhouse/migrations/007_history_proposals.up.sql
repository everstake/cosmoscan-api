CREATE TABLE IF NOT EXISTS history_proposals
(
    hpr_id           UInt64,
    hpr_tx_hash      String,
    hpr_title        String,
    hpr_description  String,
    hpr_recipient    String,
    hpr_amount       Decimal128(18),
    hpr_init_deposit Decimal128(18),
    hpr_proposer     String,
    hpr_created_at   DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMMDD(hpr_created_at)
      ORDER BY (hpr_id);