CREATE TABLE users (
  id VARCHAR(255) PRIMARY KEY,
  first_name VARCHAR(255),
  last_name VARCHAR(255),
  email VARCHAR(255),
  password VARCHAR(255),
  created_at timestamp
);

CREATE UNIQUE INDEX "email_index" ON users(email);
