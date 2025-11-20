CREATE TABLE companies (
  id UUID PRIMARY KEY,
  title VARCHAR NOT  NULL,
  owner_id UUID NOT NULL REFERENCES users(id),
  status status_enum NOT NULL DEFAULT 'active',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);