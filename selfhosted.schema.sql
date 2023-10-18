create table users(
  id uuid not null primary key,
  username text not null unique,
  password text not null,
  user_type text not null default 'regular'
);

create table sessions(
  id uuid not null primary key,
  user uuid not null references users(id)
);

create table if not exists databases (
  id uuid not null primary key,
  label text not null,
  database_variant text not null,
  owner uuid not null references users(id),
  url text not null
);

create table if not exists backup_schedules (
  id uuid not null primary key,
  database uuid not null references databases(id),
  owner uuid not null references users(id),
  frequency text not null,
  active boolean default true,
  created_at timestamp not null
);

create table if not exists database_backups (
  id uuid not null primary key,
  database uuid not null references databases(id),
  object text not null,
  created_at timestamp not null
);

create table if not exists backup_logs(
  id serial primary key,
  backup_schedule uuid not null references backup_schedules(id),
  created_at timestamp not null,
  status text not null
);
