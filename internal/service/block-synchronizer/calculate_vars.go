package synchronizer

import (
	"github.com/yanmoyy/tbi/internal/indexer"
)

func calculateGetTransactionsVarsList(blocks []indexer.Block, batchSize int) []indexer.GetTransactionsVars {
	var result []indexer.GetTransactionsVars
	if len(blocks) == 0 {
		return result
	}

	stHeight := int(blocks[0].Height)
	endHeight := int(blocks[len(blocks)-1].Height)
	endIndex := int(blocks[0].NumTxs)

	curTxs := 0

	for _, block := range blocks {
		height := int(block.Height)
		numTxs := int(block.NumTxs)

		curTxs += numTxs

		if curTxs >= batchSize {
			if height != stHeight { // need to fetch until before
				result = append(result, indexer.GetTransactionsVars{
					StartHeight: stHeight,
					EndHeight:   height,
					StartIndex:  0,
					EndIndex:    endIndex,
				})
			}
			// current block
			curTxs = 0
			stHeight = height

			if numTxs >= batchSize {
				stIndex := 0
				for ; stIndex+batchSize < numTxs; stIndex += batchSize {
					result = append(result, indexer.GetTransactionsVars{
						StartHeight: height,
						EndHeight:   height + 1,
						StartIndex:  stIndex,
						EndIndex:    stIndex + batchSize,
					})
				}
				// remaining
				result = append(result, indexer.GetTransactionsVars{
					StartHeight: height,
					EndHeight:   height + 1,
					StartIndex:  stIndex,
					EndIndex:    numTxs,
				})
				stHeight = height + 1
			}
		}
		endIndex = max(endIndex, numTxs)
	}

	if curTxs > 0 {
		result = append(result, indexer.GetTransactionsVars{
			StartHeight: stHeight,
			EndHeight:   endHeight + 1,
			StartIndex:  0,
			EndIndex:    endIndex,
		})
	}

	return result
}
