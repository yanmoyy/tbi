package database

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yanmoyy/tbi/internal/models"
)

func TestGetLastBlock(t *testing.T) {
	c := getTestClient(t)

	type input struct {
		blocks []models.Block
	}
	type expected struct {
		height   int
		totalTxs int
	}
	tester := func(t *testing.T, i input, e expected) {
		err := c.ClearAll()
		require.NoError(t, err)
		if len(i.blocks) != 0 {
			err = c.CreateBlocks(i.blocks)
			require.NoError(t, err)
		}
		height, totalTxs, err := c.GetLastBlockInfo()
		require.NoError(t, err)
		require.Equal(t, e.height, height)
		require.Equal(t, e.totalTxs, totalTxs)
	}

	t.Run("normal case", func(t *testing.T) {
		tester(t, input{
			blocks: []models.Block{
				{
					Hash:   "hash1",
					Height: 0,
				},
				{
					Hash:     "hash2",
					Height:   1,
					TotalTxs: 1,
				},
			},
		}, expected{
			height:   1,
			totalTxs: 1,
		})
	})
	t.Run("block missed", func(t *testing.T) {
		tester(t, input{
			blocks: []models.Block{
				{
					Hash:   "hash1",
					Height: 0,
				},
				{
					Hash:     "hash2",
					Height:   2,
					TotalTxs: 1,
				},
			},
		}, expected{
			height:   2, // NOTE: In this version, we don't know data is corrupted or not.
			totalTxs: 1,
		})
	})
	t.Run("block empty", func(t *testing.T) {
		tester(t, input{
			blocks: []models.Block{},
		}, expected{
			height:   -1,
			totalTxs: 0,
		})
	})
}

// Deprecated: ignore data corruption
//
// func TestGetLastHeight(t *testing.T) {
// 	c := getTestClient(t)
//
// 	type input struct {
// 		blocks []models.Block
// 	}
// 	type expected struct {
// 		height int32
// 		missed bool
// 	}
//
// 	tester := func(t *testing.T, i input, e expected) {
// 		err := c.ClearAll()
// 		require.NoError(t, err)
// 		if len(i.blocks) != 0 {
// 			err = c.CreateBlocks(i.blocks)
// 			require.NoError(t, err)
// 		}
// 		height, missed, err := c.GetLastHeight()
// 		require.NoError(t, err)
// 		require.Equal(t, e.height, height)
// 		require.Equal(t, e.missed, missed)
// 	}
//
// 	t.Run("normal case", func(t *testing.T) {
// 		tester(t, input{
// 			blocks: []models.Block{
// 				{
// 					Hash:   "hash1",
// 					Height: 0,
// 				},
// 				{
// 					Hash:   "hash2",
// 					Height: 1,
// 				},
// 			},
// 		}, expected{
// 			height: 1,
// 			missed: false,
// 		})
// 	})
// 	t.Run("block missed", func(t *testing.T) {
// 		tester(t, input{
// 			blocks: []models.Block{
// 				{
// 					Hash:   "hash1",
// 					Height: 0,
// 				},
// 				{
// 					Hash:   "hash2",
// 					Height: 2,
// 				},
// 			},
// 		}, expected{
// 			height: 2,
// 			missed: true,
// 		})
// 	})
// 	t.Run("block empty", func(t *testing.T) {
// 		tester(t, input{
// 			blocks: []models.Block{},
// 		}, expected{
// 			height: -1,
// 			missed: false,
// 		})
// 	})
// }
//
// func TestGetMissingHeight(t *testing.T) {
// 	c := getTestClient(t)
//
// 	type input struct {
// 		blocks []models.Block
// 	}
// 	type expected struct {
// 		missing []int32
// 	}
//
// 	tester := func(t *testing.T, i input, e expected) {
// 		err := c.ClearAll()
// 		require.NoError(t, err)
// 		err = c.CreateBlocks(i.blocks)
// 		require.NoError(t, err)
//
// 		missing, err := c.GetMissingHeights()
// 		require.NoError(t, err)
// 		require.Equal(t, e.missing, missing)
// 	}
//
// 	t.Run("no missing", func(t *testing.T) {
// 		tester(t, input{
// 			[]models.Block{
// 				{
// 					Hash:   "hash1",
// 					Height: 0,
// 				},
// 				{
// 					Hash:   "hash2",
// 					Height: 1,
// 				},
// 			},
// 		}, expected{
// 			missing: nil,
// 		})
// 	})
// 	t.Run("missing one", func(t *testing.T) {
// 		tester(t, input{
// 			[]models.Block{
// 				{
// 					Hash:   "hash1",
// 					Height: 0,
// 				},
// 				{
// 					Hash:   "hash2",
// 					Height: 2,
// 				},
// 			},
// 		}, expected{
// 			missing: []int32{1},
// 		})
// 	})
// 	t.Run("missing multiple", func(t *testing.T) {
// 		tester(t, input{
// 			[]models.Block{
// 				{
// 					Hash:   "hash1",
// 					Height: 0,
// 				},
// 				{
// 					Hash:   "hash2",
// 					Height: 5,
// 				},
// 			},
// 		}, expected{
// 			missing: []int32{1, 2, 3, 4},
// 		})
// 	})
// }
