package database

import (
	"fmt"

	"github.com/yanmoyy/tbi/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client struct {
	db *gorm.DB
}

func NewClient(cfg config.DB) *Client {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("db.Connect: " + err.Error())
	}
	return &Client{db: db}
}

func (c *Client) ClearAll() error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		err := c.ClearBlocks()
		if err != nil {
			return err
		}
		err = c.ClearTransactions()
		if err != nil {
			return err
		}
		err = c.ClearTokenBalances()
		if err != nil {
			return err
		}
		err = c.ClearTokenTransfers()
		if err != nil {
			return err
		}
		return nil
	})
}
