package models

import (
	"encoding/json"
	"time"
)

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

type GnoEvent struct {
	Type    string         `json:"type"`
	Func    string         `json:"func"`
	PkgPath string         `json:"pkg_path"`
	Attrs   []GnoEventAttr `json:"attrs"`
}

type GnoEventAttr struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
