package synchronizer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/yanmoyy/tbi/internal/database"
	"github.com/yanmoyy/tbi/internal/indexer"
)

const (
	blockBatchSize  = 2000
	txBatchSize     = 20
	backfillTimeout = time.Second * 60
)

type Service struct {
	indexer *indexer.Client
	db      *database.Client
	sqs     *sqs.Client
}

func New(client *indexer.Client, db *database.Client, sqsClient *sqs.Client) *Service {
	return &Service{
		indexer: client,
		db:      db,
		sqs:     sqsClient,
	}
}

func (s *Service) RunBackfill() {
	lastHeight, totalTxs, err := s.db.GetLastBlockInfo()
	if err != nil {
		slog.Error("GetLastHeight", "err", err)
	}
	slog.Info("Last height", "height", lastHeight)
	// start to fetch block
	ctx, cancel := context.WithTimeout(context.Background(), backfillTimeout)
	defer cancel()

	err = s.tryFetchAll(ctx, lastHeight, totalTxs)
	if err != nil {
		slog.Error("tryFetchAll", "err", err)
	}
}

type counter struct {
	Blocks int
	Txs    int
}

func (s *Service) tryFetchAll(ctx context.Context, lastHeight, totalTxs int) error {
	var st, en int
	st = lastHeight + 1
	en = st + blockBatchSize

	counter := counter{
		Blocks: 0,
		Txs:    0,
	}

	for {
		resp, err := s.indexer.GetBlocks(ctx, st, en)
		if err != nil {
			return fmt.Errorf("GetBlocks: %w", err)
		}
		blockCount := len(resp.Blocks)
		slog.Info("Blocks fetched", "st", st, "en", en, "count", blockCount)

		counter.Blocks += blockCount

		if blockCount == 0 { // no new blocks
			break
		}

		err = s.db.CreateBlocks(resp.ToModel())
		if err != nil {
			return fmt.Errorf("CreateBlocks: %w", err)
		}

		txsCount := int(resp.Blocks[blockCount-1].TotalTxs) - totalTxs

		if txsCount > 0 {
			totalTxs += txsCount

			overflow := txsCount - txBatchSize
			if overflow > 0 {
				blocks := resp.Blocks
				if st == 0 { // NOTE: avoid first block
					blocks = blocks[1:]
				}
				varsList := calculateGetTransactionsVarsList(blocks, txBatchSize)
				for _, vars := range varsList {
					n, err := s.fetchAndSaveTransactions(ctx,
						vars.StartHeight,
						vars.EndHeight,
						vars.StartIndex,
						vars.EndIndex,
					)
					if err != nil {
						return fmt.Errorf("fetchAndSaveTransactions: %w", err)
					}
					counter.Txs += n
				}
			} else {
				if st == 0 { // NOTE: avoid first block
					st = 1
				}
				n, err := s.fetchAndSaveTransactions(ctx, st, en, 0, txsCount)
				if err != nil {
					return fmt.Errorf("fetchAndSaveTransactions: %w", err)
				}
				counter.Txs += n
			}
		}
		if blockCount < blockBatchSize { // reach end
			break
		}
		st = en
		en += blockBatchSize
	}

	slog.Info("Synchronized blocks", "blocks", counter.Blocks, "txs", counter.Txs)
	return nil
}

func (s *Service) fetchAndSaveTransactions(ctx context.Context,
	st, en, startIndex, endIndex int) (int, error) {
	resp, err := s.indexer.GetTransactions(ctx, indexer.GetTransactionsVars{
		StartHeight: st,
		EndHeight:   en,
		StartIndex:  startIndex,
		EndIndex:    endIndex,
	})
	if err != nil {
		return 0, err
	}
	n := len(resp.Transactions)
	slog.Info("Transactions fetched", "st", st, "en", en, "count", n)
	if n == 0 {
		return 0, nil
	}
	transactions, err := resp.ToModel()
	if err != nil {
		return n, err
	}
	err = s.db.CreateTransactions(transactions)
	if err != nil {
		return n, err
	}
	return n, nil
}
