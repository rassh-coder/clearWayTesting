create table if not exists sessions (
    id text primary key default encode(gen_random_bytes(16),'hex'),
    uid bigint not null, -- user id
    created_at timestamptz not null default now(), -- Не добавил доп поле, потому что created_at = дате входа пользователя
    ip text not null,
    is_active bool not null -- Активная сессия
);