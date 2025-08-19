CREATE TABLE blocks (
  hash VARCHAR(64) PRIMARY KEY,
  height INT UNIQUE NOT NULL,
  time TIMESTAMP NOT NULL,
  total_txs INT NOT NULL,
  num_txs INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  index INT NOT NULL,
  hash VARCHAR(64) NOT NULL,
  success BOOLEAN NOT NULL,
  block_height INT NOT NULL,
  gas_wanted INT NOT NULL,
  gas_used BIGINT NOT NULL,
  memo TEXT NOT NULL,
  gas_fee JSONB,
  messages JSONB NOT NULL,
  response JSONB NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  FOREIGN KEY (block_height) REFERENCES blocks (height),
  CONSTRAINT unique_index_block_height UNIQUE (index, block_height)
);

CREATE INDEX idx_transactions_block_height ON transactions (block_height);

CREATE TABLE token_balances (
  address VARCHAR(90) NOT NULL,
  token_path VARCHAR(255) NOT NULL,
  amount BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (address, token_path),
  CHECK (amount >= 0)
);

CREATE INDEX idx_token_balances_address ON token_balances (address);

CREATE INDEX idx_token_balances_token_path ON token_balances (token_path);

CREATE TABLE token_transfers (
  id SERIAL PRIMARY KEY,
  from_address VARCHAR(90) NOT NULL,
  to_address VARCHAR(90) NOT NULL,
  token_path VARCHAR(255) NOT NULL,
  amount BIGINT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  CHECK (amount >= 0)
);

CREATE INDEX idx_token_transfers_from_address ON token_transfers (from_address);

CREATE INDEX idx_token_transfers_to_address ON token_transfers (to_address);

