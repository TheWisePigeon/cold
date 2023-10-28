CREATE TABLE users (
  id uuid not null primary key,
  username text not null unique,
  password text not null
);

CREATE TABLE integrations (
  id uuid not null primary key,
  name text not null
);

CREATE TABLE user_integrations (
  id uuid not null primary key,
  integration uuid not null,
  owner uuid not null,
  FOREIGN KEY (integration) REFERENCES integrations(id),
  FOREIGN KEY (owner) REFERENCES users(id)
);

CREATE TABLE integration_credentials (
  id uuid not null primary key,
  integration uuid not null,
  key text not null,
  value text not null,
  FOREIGN KEY (integration) REFERENCES integrations(id)
);

CREATE TABLE sessions(
  id uuid not null primary key,
  session_user uuid not null,
  FOREIGN KEY (session_user) REFERENCES users(id)
);
