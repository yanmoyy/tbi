package indexer

import (
	"strconv"
	"testing"
	"time"

	"github.com/yanmoyy/tbi/internal/models"
	// require
	"github.com/stretchr/testify/require"
)

func TestConvert(t *testing.T) {
	t.Run("convert slice", func(t *testing.T) {
		slice := []int{1, 2, 3}
		got, err := convert(slice, func(i int) (string, error) {
			return strconv.Itoa(i), nil
		})
		if err != nil {
			t.Fatal(err)
		}
		want := []string{"1", "2", "3"}
		require.Equal(t, got, want)
	})

	t.Run("convert block", func(t *testing.T) {
		now := time.Now()
		blocks := []block{
			{
				Hash:     "hash",
				Height:   1,
				Time:     now,
				NumTxs:   2,
				TotalTxs: 3,
			},
		}
		got, err := convert(blocks, blockConvertor)
		if err != nil {
			t.Fatal(err)
		}
		want := []models.Block{
			{
				Hash:     "hash",
				Height:   1,
				Time:     now,
				NumTxs:   2,
				TotalTxs: 3,
			},
		}
		require.Equal(t, got, want)
	})

	t.Run("convert transaction", func(t *testing.T) {

		transactions := []transaction{
			{
				Index:       1,
				Hash:        "hash",
				Success:     true,
				BlockHeight: 2,
				GasWanted:   3,
				GasUsed:     4,
				Memo:        "memo",
				GasFee: gasFee{
					Amount: 5,
					Denom:  "denom",
				},
				Messages: []transactionMessage{
					{
						Route:   "route",
						TypeURL: "send", // NOTE: send for bankMsgSend
						Value: messageValue{
							BankMsgSend: bankMsgSend{
								FromAddress: "fromAddress",
								ToAddress:   "toAddress",
								Amount:      "amount",
							},
						},
					},
				},
				Response: transactionResponse{
					Log:   "log",
					Info:  "info",
					Error: "error",
					Data:  "data",
					Events: []responseEvent{
						{
							GNOEvent: gnoEvent{
								Type:    "type",
								Func:    "func",
								PkgPath: "pkgPath",
								Attrs: []gnoEventAttr{
									{
										Key:   "key",
										Value: "value",
									},
								},
							},
						},
					},
				},
			},
		}
		got, err := convert(transactions, transactionConvertor)
		if err != nil {
			t.Fatal(err)
		}
		want := []models.Transaction{
			{
				Index:       1,
				Hash:        "hash",
				Success:     true,
				BlockHeight: 2,
				GasWanted:   3,
				GasUsed:     4,
				Memo:        "memo",
				GasFee: models.Coin{
					Amount: 5,
					Denom:  "denom",
				},
				Messages: []models.TransactionMessage{
					{
						Route:   "route",
						TypeURL: "send",
						Value: []byte(
							`{"from_address":"fromAddress","to_address":"toAddress","amount":"amount"}`,
						),
					},
				},
				Response: models.TransactionResponse{
					Log:   "log",
					Info:  "info",
					Error: "error",
					Data:  "data",
					Events: []models.GnoEvent{
						{
							Type:    "type",
							Func:    "func",
							PkgPath: "pkgPath",
							Attrs: []models.GnoEventAttr{
								{
									Key:   "key",
									Value: "value",
								},
							},
						},
					},
				},
			},
		}
		require.Equal(t, got, want)
	})
}
