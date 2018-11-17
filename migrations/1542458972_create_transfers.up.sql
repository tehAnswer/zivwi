CREATE TABLE transfers (
  id VARCHAR(255) PRIMARY KEY,
  from_account_id VARCHAR(255) NOT NULL,
  to_account_id VARCHAR(255) NOT NULL,
  amount NUMERIC NOT NULL,
  message VARCHAR(255),
  status VARCHAR(255),
  error VARCHAR(255),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL
);
