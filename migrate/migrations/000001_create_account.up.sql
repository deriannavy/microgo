CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS account(
   id bigserial PRIMARY KEY,
   username  citext UNIQUE NOT NULL,
   email     citext UNIQUE NOT NULL,
   password  bytea,
   created_At TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
   is_active BOOLEAN NOT NULL DEFAULT FALSE
)

