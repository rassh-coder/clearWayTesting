create table if not exists assets (
    name text not null,
    uid bigint not null, -- user id
    data bytea not null,
    created_at timestamptz not null default now(),
    primary key (name, uid)
);