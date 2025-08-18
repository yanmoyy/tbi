package models

import (
	"encoding/json"
	"time"
)

// Database models

type Block struct {
	Hash      string    `gorm:"primaryKey;type:varchar(64)"`
	Height    int       `gorm:"unique;not null;index"`
	Time      time.Time `gorm:"not null"`
	NumTxs    int       `gorm:"not null"`
	TotalTxs  int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Transaction struct {
	ID          uint            `gorm:"primaryKey;autoIncrement;"`
	Index       int             `gorm:"not null"`
	Hash        string          `gorm:"type:varchar(64);not null"`
	Success     bool            `gorm:"not null"`
	BlockHeight int             `gorm:"not null;index"`
	GasWanted   int             `gorm:"not null"`
	GasUsed     int             `gorm:"not null"`
	Memo        string          `gorm:"type:text;not null"`
	GasFee      json.RawMessage `gorm:"type:jsonb"`
	Messages    json.RawMessage `gorm:"type:jsonb;not null"`
	Response    json.RawMessage `gorm:"type:jsonb;not null"`
	CreatedAt   time.Time       `gorm:"autoCreateTime"`
}

type TokenBalance struct {
	Address   string `gorm:"primaryKey;type:varchar(90)"`
	TokenPath string `gorm:"primaryKey;type:varchar(255)"`
	Amount    int64  `gorm:"not null"`
}

type TokenTransfer struct {
	ID          uint      `gorm:"primaryKey;autoIncrement;"`
	FromAddress string    `gorm:"type:varchar(90);not null"`
	ToAddress   string    `gorm:"type:varchar(90);not null"`
	TokenPath   string    `gorm:"type:varchar(255);not null"`
	Amount      int64     `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

// Event
type EventFunc string

const (
	Mint     EventFunc = "Mint"
	Burn     EventFunc = "Burn"
	Transfer EventFunc = "Transfer"
)

type TransferEvent struct {
	Func      EventFunc `json:"func"`
	TokenPath string    `json:"token"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	Value     int64     `json:"value"`
}
