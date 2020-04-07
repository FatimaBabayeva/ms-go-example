-- +migrate Up
create table if not exists message
(
    id              bigserial       not null primary key,
    text            varchar(256)    not null,
    status          varchar(16)     not null default 'CREATED',
    created_at      timestamp       not null default now(),
    updated_at      timestamp       not null default now()
);
