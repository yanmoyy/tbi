package indexer

import (
	"encoding/json"
	"time"

	"github.com/yanmoyy/tbi/internal/models"
)

type Block struct {
	Hash     string    `graphql:"hash"`
	Height   int32     `graphql:"height"`
	Time     time.Time `graphql:"time"`
	NumTxs   int32     `graphql:"num_txs"`
	TotalTxs int32     `graphql:"total_txs"`
}

func (b *Block) ToModel() models.Block {
	return models.Block{
		Hash:     b.Hash,
		Height:   int(b.Height),
		Time:     b.Time,
		NumTxs:   int(b.NumTxs),
		TotalTxs: int(b.TotalTxs),
	}
}

type Transaction struct {
	Index       int32                `graphql:"index"`
	Hash        string               `graphql:"hash"`
	Success     bool                 `graphql:"success"`
	BlockHeight int32                `graphql:"block_height"`
	GasWanted   int32                `graphql:"gas_wanted"`
	GasUsed     int32                `graphql:"gas_used"`
	Memo        string               `graphql:"memo"`
	GasFee      gasFee               `graphql:"gas_fee"`
	Messages    []transactionMessage `graphql:"messages"`
	Response    transactionResponse  `graphql:"response"`
}

// error reason: need to use json.Marshal
func (t *Transaction) ToModel() (models.Transaction, error) {
	gasFeeJson, err := json.Marshal(t.GasFee)
	if err != nil {
		return models.Transaction{}, err
	}
	messagesJson, err := json.Marshal(t.Messages)
	if err != nil {
		return models.Transaction{}, err
	}
	responseJson, err := json.Marshal(t.Response)
	if err != nil {
		return models.Transaction{}, err
	}
	return models.Transaction{
		Index:       int(t.Index),
		Hash:        t.Hash,
		Success:     t.Success,
		BlockHeight: int(t.BlockHeight),
		GasWanted:   int(t.GasWanted),
		GasUsed:     int(t.GasUsed),
		Memo:        t.Memo,
		GasFee:      gasFeeJson,
		Messages:    messagesJson,
		Response:    responseJson,
	}, nil
}

type gasFee struct {
	Amount int32  `graphql:"amount"`
	Denom  string `graphql:"denom"`
}

type transactionMessage struct {
	Route   string       `graphql:"route"`
	TypeURL string       `graphql:"typeUrl"`
	Value   messageValue `graphql:"value"`
}

type messageValue struct {
	BankMsgSend   bankMsgSend   `graphql:"... on BankMsgSend"`
	MsgAddPackage msgAddPackage `graphql:"... on MsgAddPackage"`
	MsgCall       msgCall       `graphql:"... on MsgCall"`
	MsgRun        msgRun        `graphql:"... on MsgRun"`
}

type transactionResponse struct {
	Log    string          `graphql:"log"`
	Info   string          `graphql:"info"`
	Error  string          `graphql:"error"`
	Data   string          `graphql:"data"`
	Events []responseEvent `graphql:"events"`
}

type responseEvent struct {
	GNOEvent gnoEvent `graphql:"... on GnoEvent"`
}

// message values

type bankMsgSend struct {
	FromAddress string `graphql:"from_address" json:"from_address"`
	ToAddress   string `graphql:"to_address" json:"to_address"`
	Amount      string `graphql:"amount" json:"amount"`
}

type msgAddPackage struct {
	Creator string      `graphql:"creator" json:"creator"`
	Send    string      `graphql:"send" json:"send"`
	Package packageData `graphql:"package" json:"package"`
}

type msgCall struct {
	PkgPath string   `graphql:"pkg_path" json:"pkg_path"`
	Func    string   `graphql:"func" json:"func"`
	Send    string   `graphql:"send" json:"send"`
	Caller  string   `graphql:"caller" json:"caller"`
	Args    []string `graphql:"args" json:"args"`
}

type msgRun struct {
	Caller  string      `graphql:"caller" json:"caller"`
	Send    string      `graphql:"send" json:"send"`
	Package packageData `graphql:"package" json:"package"`
}

/// #########

type packageData struct {
	Name  string        `graphql:"name"`
	Path  string        `graphql:"path"`
	Files []packageFile `graphql:"files"`
}

type packageFile struct {
	Name string `graphql:"name"`
	Body string `graphql:"body"`
}

// response events
type gnoEvent struct {
	Type    string         `graphql:"type"`
	Func    string         `graphql:"func"`
	PkgPath string         `graphql:"pkg_path"`
	Attrs   []gnoEventAttr `graphql:"attrs"`
}

type gnoEventAttr struct {
	Key   string `graphql:"key"`
	Value string `graphql:"value"`
}
