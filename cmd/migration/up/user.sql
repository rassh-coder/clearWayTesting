create extension if not exists pgcrypto;

create table if not exists users (
    id bigserial primary key,
    login text not null unique,
    password_hash text not null,
    created_at timestamptz not null default now()
);

insert into users (login, password_hash)
values ('alice', encode(digest('secret', 'md5'),'hex'))
    on conflict do nothing;