CREATE TABLE IF NOT EXISTS transaction(
    id bigserial PRIMARY KEY,
	account_id bigint NOT NULL,
	date date NOT NULL,
	amount numeric NOT NULL,
	place varchar(255) NOT NULL,
	description varchar(255) NOT NULL,
	tag varchar(255)[] NOT NULL,
    FOREIGN KEY (account_id) REFERENCES account(id)
)