-- +goose Up
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

-- +goose Down
DROP TABLE transactions;
