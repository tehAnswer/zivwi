CREATE TABLE accounts (
  id VARCHAR(255) PRIMARY KEY,
  balance NUMERIC NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL
);

ALTER TABLE users ADD COLUMN account_ids text[];
