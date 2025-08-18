package database

import (
	"context"

	"github.com/yanmoyy/tbi/internal/models"
)

func (c *Client) CreateTokenTransfer(ctx context.Context, transfer models.TokenTransfer) error {
	return c.db.WithContext(ctx).Create(&transfer).Error
}

func (c *Client) GetTokenTransferList(ctx context.Context) ([]models.TokenTransfer, error) {
	var history []models.TokenTransfer
	err := c.db.WithContext(ctx).
		Order("id DESC").
		Find(&history).
		Error
	return history, err
}

func (c *Client) GetTokenTransferListWithAddress(ctx context.Context, address string) ([]models.TokenTransfer, error) {
	var history []models.TokenTransfer
	err := c.db.WithContext(ctx).
		Where("from_address = ? OR to_address = ?", address, address).
		Order("id DESC").
		Find(&history).
		Error
	return history, err
}

func (c *Client) ClearTokenTransfers() error {
	const query = /*sql*/ `
TRUNCATE TABLE token_transfers CASCADE;
ALTER SEQUENCE token_transfers_id_seq RESTART WITH 1;
`
	return c.db.Exec(query).Error
}
