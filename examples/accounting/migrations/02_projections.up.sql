BEGIN;

CREATE TABLE IF NOT EXISTS books (
    id TEXT PRIMARY KEY,
    closed BOOL NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts (
    book_id TEXT NOT NULL REFERENCES books (id),
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    balance BIGINT NOT NULL,
    UNIQUE (book_id, "name")
);

CREATE TABLE IF NOT EXISTS transactions (
    book_id TEXT NOT NULL REFERENCES books (id),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    account_debited TEXT NOT NULL,
    account_credited TEXT NOT NULL,
    amount BIGINT NOT NULL
);

END;
