-- +migrate Up
create table parsers
(
    par_id     int auto_increment
        primary key,
    par_title  varchar(255)  not null,
    par_height int default 0 not null,
    constraint parsers_par_title_uindex
        unique (par_title)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

insert into parsers (par_id, par_title, par_height)
VALUES (1, 'hub3', 0);

create table validators
(
    val_cons_address     varchar(255)                                                                            not null
        primary key,
    val_address          varchar(255)                                                  default ''                not null,
    val_operator_address varchar(255)                                                  default ''                not null,
    val_cons_pub_key     varchar(255)                                                  default ''                not null,
    val_name             varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci default ''                not null,
    val_description      text                                                                                    not null,
    val_commission       decimal(8, 4)                                                 default 0.0000            not null,
    val_min_commission   decimal(8, 4)                                                 default 0.0000            not null,
    val_max_commission   decimal(8, 4)                                                 default 0.0000            not null,
    val_self_delegations decimal(20, 8)                                                                          not null,
    val_delegations      decimal(20, 8)                                                default 0.00000000        not null,
    val_voting_power     decimal(20, 8)                                                default 0.00000000        not null,
    val_website          varchar(255)                                                  default ''                not null,
    val_jailed           tinyint(1)                                                    default 0                 null,
    val_created_at       timestamp                                                     default CURRENT_TIMESTAMP not null,
    constraint validators_val_cons_address_uindex
        unique (val_cons_address)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

create table accounts
(
    acc_address    varchar(255)                             not null
        primary key,
    acc_balance    decimal(20, 8) default 0.00000000        not null,
    acc_created_at timestamp      default CURRENT_TIMESTAMP not null
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;


create table accounts
(
    acc_address    varchar(255)                             not null
        primary key,
    acc_balance    decimal(20, 8) default 0.00000000        not null,
    acc_created_at timestamp      default CURRENT_TIMESTAMP not null
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;


create table range_states
(
    rst_title      varchar(255)                        not null
        primary key,
    rst_value_1d   text                                null,
    rst_value_7d   text                                null,
    rst_value_30d  text                                null,
    rst_value_90d  text                                null,
    rst_updated_at timestamp default CURRENT_TIMESTAMP not null on update CURRENT_TIMESTAMP
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;


-- +migrate Down
drop table parsers;
drop table validators;
drop table accounts;
drop table range_states;
