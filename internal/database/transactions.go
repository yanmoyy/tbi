package database

import (
	"context"

	"github.com/yanmoyy/tbi/internal/models"
)

func (c *Client) CreateTransactionList(ctx context.Context, transactions []models.Transaction) error {
	return c.db.WithContext(ctx).Create(&transactions).Error
}

func (c *Client) ClearTransactions() error {
	const query = /*sql*/ `
TRUNCATE TABLE transactions CASCADE;
ALTER SEQUENCE transactions_id_seq RESTART WITH 1;
`
	return c.db.Exec(query).Error
}
