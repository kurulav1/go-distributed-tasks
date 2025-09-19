use app;
create table if not exists schema_migrations (name string primary key);
create table if not exists users (
  id uuid default gen_random_uuid() primary key,
  email string not null unique,
  password_hash bytes not null,
  created_at timestamp not null default now()
);
create table if not exists jobs (
  id uuid primary key,
  user_id uuid not null references users(id),
  type string not null,
  payload jsonb not null,
  enqueued_at timestamp not null
);
create index if not exists idx_jobs_user on jobs(user_id);
create index if not exists idx_jobs_type on jobs(type);
insert into schema_migrations(name) values ('001_users.sql') on conflict do nothing;
insert into schema_migrations(name) values ('002_jobs.sql') on conflict do nothing;
