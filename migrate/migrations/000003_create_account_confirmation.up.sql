CREATE TABLE IF NOT EXISTS account_confirmation(
    token bytea PRIMARY KEY,
    account_id bigint NOT NULL,
    expiry TIMESTAMP(0) WITH TIME ZONE NOT NULL
)