create table if not exists jobs (
  id uuid primary key,
  user_id uuid not null references users(id),
  type string not null,
  payload jsonb not null,
  enqueued_at timestamp not null
);
create index if not exists idx_jobs_user on jobs(user_id);
create index if not exists idx_jobs_type on jobs(type);
