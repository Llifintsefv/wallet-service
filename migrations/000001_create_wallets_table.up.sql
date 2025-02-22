CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    balance INTEGER NOT NULL DEFAULT 0,
    CONSTRAINT non_negative_balance CHECK (balance >= 0)
);
