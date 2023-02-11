CREATE TABLE IF NOT EXISTS historical_states
(
    his_price              Decimal(18, 8) default 0,
    his_market_cap         Decimal(18, 2) default 0,
    his_circulating_supply Decimal(18, 2) default 0,
    his_trading_volume     Decimal(18, 2) default 0,
    his_staked_ratio       Decimal(4, 2)  default 0,
    his_inflation_rate     Decimal(4, 2)  default 0,
    his_transactions_count UInt64         default 0,
    his_community_pool     Decimal(18, 2) default 0,
    his_top_20_weight      Decimal(4, 2)  default 0,
    his_created_at         DateTime
) ENGINE ReplacingMergeTree()
      PARTITION BY toYYYYMMDD(his_created_at)
      ORDER BY (his_created_at);

