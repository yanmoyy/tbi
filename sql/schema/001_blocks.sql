-- +goose Up
CREATE TABLE blocks (
  hash VARCHAR(64) PRIMARY KEY,
  height INT UNIQUE NOT NULL,
  time TIMESTAMP NOT NULL,
  total_txs INT NOT NULL,
  num_txs INT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE blocks;
