-- +migrate Up
create table parsers
(
    par_id     int auto_increment
        primary key,
    par_title  varchar(255)  not null,
    par_height int default 0 not null,
    constraint parsers_par_title_uindex
        unique (par_title)
);

insert into parsers (par_id, par_title, par_height) VALUES (1, "hub3", 0);

-- +migrate Down
drop table parsers;
