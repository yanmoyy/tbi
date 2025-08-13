package models

import (
	"strconv"
	"time"
)

type Block struct {
	Hash     string `gorm:"primaryKey"`
	Height   int64
	Time     time.Time
	TotalTxs int64
	NumTxs   int64
}

type Transaction struct {
	Hash        string `gorm:"primaryKey"`
	BlockHeight int64
	Success     bool
	GasUsed     int64
}

type Balance struct {
	Address   string `gorm:"primaryKey"`
	TokenPath string `gorm:"primaryKey"`
	Amount    int64
}

type Transfer struct {
	ID          uint `gorm:"primaryKey"`
	FromAddress string
	ToAddress   string
	TokenPath   string
	Amount      int64
	EventTime   time.Time
}

type Event struct {
	Type    string              `json:"type"`
	Func    string              `json:"func"`
	PkgPath string              `json:"pkg_path"`
	Attrs   []map[string]string `json:"attrs"`
}

func (e Event) IsValidTokenEvent() bool {
	if e.Type != "Transfer" {
		return false
	}
	attrs := make(map[string]string)
	for _, a := range e.Attrs {
		attrs[a["key"]] = a["value"]
	}
	if _, ok := attrs["from"]; !ok {
		return false
	}
	if _, ok := attrs["to"]; !ok {
		return false
	}
	if _, ok := attrs["value"]; !ok {
		return false
	}
	// Simple bech32 check (starts with 'g1' or empty)
	isBech32OrEmpty := func(s string) bool { return s == "" || (len(s) > 2 && s[:2] == "g1") }
	_, err := strconv.ParseInt(attrs["value"], 10, 64)
	if err != nil {
		return false
	}
	switch e.Func {
	case "Mint":
		return attrs["from"] == "" && isBech32OrEmpty(attrs["to"])
	case "Burn":
		return isBech32OrEmpty(attrs["from"]) && attrs["to"] == ""
	case "Transfer":
		return isBech32OrEmpty(attrs["from"]) && isBech32OrEmpty(attrs["to"])
	}
	return false
}
