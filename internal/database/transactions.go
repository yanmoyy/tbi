package database

import "github.com/yanmoyy/tbi/internal/models"

func (c *Client) CreateTransactions(transactions []models.Transaction) error {
	return c.db.Create(&transactions).Error
}
