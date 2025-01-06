CREATE TABLE IF NOT EXISTS account_confirmation(
    token bytea PRIMARY KEY,
    account_id bigint NOT NULL
)