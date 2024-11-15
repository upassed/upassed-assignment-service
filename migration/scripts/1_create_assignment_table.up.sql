create table if not exists assignment (
    id uuid primary key,
    form_id uuid not null,
    group_id uuid not null
);
