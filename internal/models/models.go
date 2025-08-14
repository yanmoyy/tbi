package models

import (
	"encoding/json"
	"time"
)

type Block struct {
	Hash      string    `gorm:"primaryKey;type:varchar(64)"`
	Height    int32     `gorm:"unique;not null;index"`
	Time      time.Time `gorm:"not null"`
	NumTxs    int32     `gorm:"not null"`
	TotalTxs  int32     `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Transaction struct {
	ID          uint                 `gorm:"primaryKey"`
	Index       int32                `gorm:"not null"`
	Hash        string               `gorm:"type:varchar(64);not null"`
	Success     bool                 `gorm:"not null"`
	BlockHeight int32                `gorm:"not null;index"`
	GasWanted   int32                `gorm:"not null"`
	GasUsed     int32                `gorm:"not null"`
	Memo        string               `gorm:"type:text;not null"`
	GasFee      Coin                 `gorm:"type:jsonb"`
	Messages    []TransactionMessage `gorm:"type:jsonb;not null"`
	Response    TransactionResponse  `gorm:"type:jsonb;not null"`
	CreatedAt   time.Time            `gorm:"autoCreateTime"`
}

type TransactionMessage struct {
	Route   string          `json:"route"`
	TypeURL string          `json:"typeUrl"`
	Value   json.RawMessage `json:"value"`
}

type TransactionResponse struct {
	Log    string     `json:"log"`
	Info   string     `json:"info"`
	Error  string     `json:"error"`
	Data   string     `json:"data"`
	Events []GnoEvent `json:"events"`
}

type Coin struct {
	Amount int32  `json:"amount"`
	Denom  string `json:"denom"`
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
