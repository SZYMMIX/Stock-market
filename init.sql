CREATE TABLE IF NOT EXISTS bank_stocks (
    stock_name VARCHAR(50) PRIMARY KEY,
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0)
);

CREATE TABLE IF NOT EXISTS wallets (
    id VARCHAR(100) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS wallet_stocks (
    wallet_id VARCHAR(100) REFERENCES wallets(id) ON DELETE CASCADE,
    stock_name VARCHAR(50),
    quantity INT NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    PRIMARY KEY (wallet_id, stock_name)
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id SERIAL PRIMARY KEY,
    type VARCHAR(10) NOT NULL,
    wallet_id VARCHAR(100) NOT NULL,
    stock_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);