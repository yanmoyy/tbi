-- +goose Up
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

-- +goose Down
DROP TABLE token_transfers;
