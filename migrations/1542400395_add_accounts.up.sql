CREATE TABLE accounts (
  id VARCHAR(255) PRIMARY KEY,
  balance NUMERIC,
  created_at TIMESTAMP
);

ALTER TABLE users ADD COLUMN accounts text[];
