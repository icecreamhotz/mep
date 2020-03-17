CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE roles AS ENUM('ADMIN', 'EMPLOYEE');

CREATE TABLE backoffice_users (
  id uuid DEFAULT uuid_generate_v4 (),
  username VARCHAR(50) UNIQUE,
  password CHAR(60),
  email VARCHAR(100) UNIQUE,
  name VARCHAR(100),
  lastname VARCHAR(100),
  role roles default 'EMPLOYEE',
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP
);