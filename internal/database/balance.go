package database

import (
	"context"
	"fmt"

	"github.com/yanmoyy/tbi/internal/models"
	"gorm.io/gorm"
)

func (c *Client) UpdateTokenBalance(ctx context.Context, address string, tokenPath string, amount int64, increment bool) error {
	balance, err := c.GetTokenBalance(ctx, address, tokenPath)

	if err == gorm.ErrRecordNotFound {
		if increment {
			return c.CreateTokenBalance(ctx, address, tokenPath, amount)
		} else {
			return fmt.Errorf("no balance to decrease")
		}
	}
	if err != nil {
		return err
	}

	var newAmount int64
	if increment {
		newAmount = balance.Amount + amount
	} else {
		newAmount = balance.Amount - amount
		if newAmount < 0 {
			return fmt.Errorf("insufficient balance")
		}
	}
	return c.db.WithContext(ctx).
		Save(&models.TokenBalance{
			Address:   address,
			TokenPath: tokenPath,
			Amount:    newAmount,
		}).Error
}

func (c *Client) TransferTokenBalance(ctx context.Context, from, to string, tokenPath string, amount int64) error {
	return c.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := c.UpdateTokenBalance(ctx, from, tokenPath, amount, false)
		if err != nil {
			return err
		}
		return c.UpdateTokenBalance(ctx, to, tokenPath, amount, true)
	})
}

func (c *Client) CreateTokenBalance(ctx context.Context, address string, tokenPath string, amount int64) error {
	return c.db.WithContext(ctx).
		Create(&models.TokenBalance{
			Address:   address,
			TokenPath: tokenPath,
			Amount:    amount,
		}).Error
}

func (c *Client) GetTokenBalance(ctx context.Context, address string, tokenPath string) (models.TokenBalance, error) {
	var balance models.TokenBalance
	err := c.db.WithContext(ctx).
		Where("address = ? AND token_path = ?", address, tokenPath).
		First(&balance).
		Error
	return balance, err
}

func (c *Client) GetTokenBalanceList(ctx context.Context) ([]models.TokenBalance, error) {
	var balances []models.TokenBalance
	err := c.db.WithContext(ctx).
		Order("address ASC, token_path ASC").
		Find(&balances).
		Error
	return balances, err
}

func (c *Client) GetTokenBalanceListWithAddress(ctx context.Context, address string) ([]models.TokenBalance, error) {
	var balances []models.TokenBalance
	err := c.db.WithContext(ctx).
		Where("address = ?", address).
		Order("token_path ASC").
		Find(&balances).
		Error
	return balances, err
}

func (c *Client) GetTokenBalanceListWithToken(ctx context.Context, tokenPath string) ([]models.TokenBalance, error) {
	var balances []models.TokenBalance
	err := c.db.WithContext(ctx).
		Where("token_path = ?", tokenPath).
		Order("address ASC").
		Find(&balances).
		Error
	return balances, err
}

func (c *Client) ClearTokenBalances() error {
	return c.db.Exec("TRUNCATE TABLE token_balances CASCADE").Error
}
