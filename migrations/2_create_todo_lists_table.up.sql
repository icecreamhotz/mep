CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE todo_lists (
  id uuid DEFAULT uuid_generate_v4 (),
  name VARCHAR(255) NOT NULL,
  done boolean DEFAULT FALSE NOT NULL,
  image VARCHAR(255),
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
)