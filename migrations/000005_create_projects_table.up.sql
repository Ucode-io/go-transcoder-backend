CREATE TABLE projects (
  id UUID PRIMARY KEY,
  title VARCHAR NOT NULL,
  access_key VARCHAR NOT NULL,
  secret_key VARCHAR NOT NULL,
  company_id UUID NOT NULL REFERENCES companies(id),
  owner_id UUID NOT NULL UNIQUE REFERENCES users(id),
  status status_enum NOT NULL,
  storage_id UUID NOT NULL REFERENCES storages(id),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);
