package indexer

import (
	"time"
)

type block struct {
	Hash     string    `graphql:"hash"`
	Height   int32     `graphql:"height"`
	Time     time.Time `graphql:"time"`
	NumTxs   int32     `graphql:"num_txs"`
	TotalTxs int32     `graphql:"total_txs"`
}

type transaction struct {
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
