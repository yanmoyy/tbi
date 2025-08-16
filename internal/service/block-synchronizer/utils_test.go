package synchronizer

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yanmoyy/tbi/internal/indexer"
)

func TestCalculateGetTransactionsVarsList(t *testing.T) {
	type input struct {
		blocks []indexer.Block
		batch  int
	}

	type expected struct {
		vars []indexer.GetTransactionsVars
	}

	tester := func(t *testing.T, input input, expected expected) {
		actual := calculateGetTransactionsVarsList(input.blocks, input.batch)
		require.Equal(t, expected.vars, actual)
	}

	t.Run("batch size is smaller than txs number", func(t *testing.T) {
		tester(t, input{
			blocks: []indexer.Block{
				{
					Height: 1,
					NumTxs: 5,
				},
				{
					Height: 2,
					NumTxs: 5,
				},
			},
			batch: 3,
		}, expected{
			vars: []indexer.GetTransactionsVars{
				{
					StartHeight: 1,
					EndHeight:   2,
					StartIndex:  0,
					EndIndex:    3,
				},
				{
					StartHeight: 1,
					EndHeight:   2,
					StartIndex:  3,
					EndIndex:    5,
				},
				{
					StartHeight: 2,
					EndHeight:   3,
					StartIndex:  0,
					EndIndex:    3,
				},
				{
					StartHeight: 2,
					EndHeight:   3,
					StartIndex:  3,
					EndIndex:    5,
				},
			},
		})
	})

	t.Run("batch size is equal to txs number", func(t *testing.T) {
		tester(t, input{
			blocks: []indexer.Block{
				{
					Height: 1,
					NumTxs: 5,
				},
				{
					Height: 2,
					NumTxs: 5,
				},
			},
			batch: 5,
		}, expected{
			vars: []indexer.GetTransactionsVars{
				{
					StartHeight: 1,
					EndHeight:   2,
					StartIndex:  0,
					EndIndex:    5,
				},
				{
					StartHeight: 2,
					EndHeight:   3,
					StartIndex:  0,
					EndIndex:    5,
				},
			},
		})
	})

	t.Run("different heights", func(t *testing.T) {
		tester(t, input{
			blocks: []indexer.Block{
				{
					Height: 1,
					NumTxs: 2,
				},
				{
					Height: 2,
					NumTxs: 3,
				},
				{
					Height: 3,
					NumTxs: 4,
				},
				{
					Height: 4,
					NumTxs: 5,
				},
			},
			batch: 4,
		}, expected{
			vars: []indexer.GetTransactionsVars{
				{
					StartHeight: 1,
					EndHeight:   2,
					StartIndex:  0,
					EndIndex:    2,
				},
				{
					StartHeight: 2,
					EndHeight:   3,
					StartIndex:  0,
					EndIndex:    3,
				},
				{
					StartHeight: 3,
					EndHeight:   4,
					StartIndex:  0,
					EndIndex:    4,
				},
				{
					StartHeight: 4,
					EndHeight:   5,
					StartIndex:  0,
					EndIndex:    4,
				},
				{
					StartHeight: 4,
					EndHeight:   5,
					StartIndex:  4,
					EndIndex:    5,
				},
			},
		})
	})

	t.Run("batch size is bigger than total txs", func(t *testing.T) {
		tester(t, input{
			blocks: []indexer.Block{
				{
					Height: 1,
					NumTxs: 5,
				},
				{
					Height: 2,
					NumTxs: 5,
				},
				{
					Height: 3,
					NumTxs: 5,
				},
				{
					Height: 4,
					NumTxs: 5,
				},
				{
					Height: 5,
					NumTxs: 6,
				},
			},
			batch: 100,
		}, expected{
			vars: []indexer.GetTransactionsVars{
				{
					StartHeight: 1,
					EndHeight:   6,
					StartIndex:  0,
					EndIndex:    6,
				},
			},
		})
	})

	t.Run("realistic case", func(t *testing.T) {
		tester(t, input{
			blocks: []indexer.Block{
				{
					Height: 0,
					NumTxs: 36,
				},
				{
					Height: 1,
					NumTxs: 0,
				},
				{
					Height: 2,
					NumTxs: 0,
				},
				{
					Height: 3,
					NumTxs: 0,
				},
			},
			batch: 10,
		}, expected{
			vars: []indexer.GetTransactionsVars{
				{
					StartHeight: 0,
					EndHeight:   1,
					StartIndex:  0,
					EndIndex:    10,
				},
				{
					StartHeight: 0,
					EndHeight:   1,
					StartIndex:  10,
					EndIndex:    20,
				},
				{
					StartHeight: 0,
					EndHeight:   1,
					StartIndex:  20,
					EndIndex:    30,
				},
				{
					StartHeight: 0,
					EndHeight:   1,
					StartIndex:  30,
					EndIndex:    36,
				},
			},
		})
	})
}
