package graphql

import (
	"time"
)

type block struct {
	Hash     string    `graphql:"hash"`
	Height   int64     `graphql:"height"`
	Time     time.Time `graphql:"time"`
	NumTxs   int64     `graphql:"num_txs"`
	TotalTxs int64     `graphql:"total_txs"`
}

type transaction struct {
	Index       int64               `graphql:"index"`
	Hash        string              `graphql:"hash"`
	Success     bool                `graphql:"success"`
	BlockHeight int64               `graphql:"block_height"`
	GasWanted   int64               `graphql:"gas_wanted"`
	GasUsed     int64               `graphql:"gas_used"`
	Memo        string              `graphql:"memo"`
	GasFee      gasFee              `graphql:"gas_fee"`
	Messages    []message           `graphql:"messages"`
	Response    transactionResponse `graphql:"response"`
}

type gasFee struct {
	Amount int64  `graphql:"amount"`
	Denom  string `graphql:"denom"`
}

type message struct {
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
	FromAddress string `graphql:"from_address"`
	ToAddress   string `graphql:"to_address"`
	Amount      string `graphql:"amount"`
}

type msgAddPackage struct {
	Creator string      `graphql:"creator"`
	Deposit string      `graphql:"deposit"`
	Package packageData `graphql:"package"`
}

type msgCall struct {
	PkgPath string   `graphql:"pkg_path"`
	Func    string   `graphql:"func"`
	Send    string   `graphql:"send"`
	Caller  string   `graphql:"caller"`
	Args    []string `graphql:"args"`
}

type msgRun struct {
	Caller  string      `graphql:"caller"`
	Send    string      `graphql:"send"`
	Package packageData `graphql:"package"`
}

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
	Type    string      `graphql:"type"`
	Func    string      `graphql:"func"`
	PkgPath string      `graphql:"pkg_path"`
	Attrs   []eventAttr `graphql:"attrs"`
}

type eventAttr struct {
	Key   string `graphql:"key"`
	Value string `graphql:"value"`
}
