CREATE TABLE wallets (
    id UUID PRIMARY KEY,
    balance DECIMAL NOT NULL DEFAULT 0,
    CONSTRAINT non_negative_balance CHECK (balance >= 0)
);
