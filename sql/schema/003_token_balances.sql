-- +goose Up
CREATE TABLE token_balances (
  address VARCHAR(90) NOT NULL,
  token_path VARCHAR(255) NOT NULL,
  amount BIGINT NOT NULL DEFAULT 0,
  PRIMARY KEY (address, token_path),
  CHECK (amount >= 0)
);

CREATE INDEX idx_token_balances_address ON token_balances (address);

CREATE INDEX idx_token_balances_token_path ON token_balances (token_path);

-- +goose Down
DROP TABLE token_balances;
