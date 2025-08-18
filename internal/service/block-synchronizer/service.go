package synchronizer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/yanmoyy/tbi/internal/database"
	"github.com/yanmoyy/tbi/internal/indexer"
	"github.com/yanmoyy/tbi/internal/sqs"
)

const (
	blockBatchSize   = 2000
	txBatchSize      = 20
	backfillTimeout  = time.Second * 60
	maxSubscribeWait = time.Second * 5
)

type Service struct {
	indexer *indexer.Client
	db      *database.Client
	sqs     *sqs.Client
}

func New(indexer *indexer.Client, db *database.Client, producer *sqs.Client) *Service {
	return &Service{
		indexer: indexer,
		db:      db,
		sqs:     producer,
	}
}

func (s *Service) RunBackFill(ctx context.Context) error {
	lastHeight, lastTotalTxs, err := s.db.GetLastBlockInfo(ctx)
	if err != nil {
		return fmt.Errorf("GetLastHeight: %w", err)
	}

	firstHeight, err := s.startSubscription(ctx)
	if err != nil {
		return fmt.Errorf("startSubscription: %w", err)
	}

	// start to backfill
	ctxB, cancel := context.WithTimeout(ctx, backfillTimeout)
	defer cancel()

	err = s.tryFetchAll(ctxB, lastHeight+1, firstHeight, lastTotalTxs)
	if err != nil {
		return fmt.Errorf("tryFetchAll: %w", err)
	}
	return nil
}

func (s *Service) startSubscription(ctx context.Context) (int, error) {
	blockCh := make(chan indexer.Block)
	done := make(chan struct{})

	ctxS, cancel := context.WithCancel(ctx)
	s.indexer.SubscribeBlocks(ctxS, blockCh, done)
	heightChan := make(chan int)
	initialized := false

	go func() {
		for {
			select {
			case <-done:
				return
			case block := <-blockCh:
				slog.Info("Block received", "height", block.Height)
				if !initialized {
					heightChan <- int(block.Height)
					initialized = true
				}
				err := s.db.CreateBlock(ctx, block.ToModel())
				if err != nil {
					slog.Error("CreateBlock", "err", err)
					cancel()
					return
				}
				if block.NumTxs > 0 {
					// process transactions event
					varsList := calculateGetTransactionsVarsList([]indexer.Block{block}, txBatchSize)
					for _, vars := range varsList {
						_, err := s.syncTransactions(ctx,
							vars.StartHeight,
							vars.EndHeight,
							vars.StartIndex,
							vars.EndIndex,
						)
						if err != nil {
							slog.Error("syncTransactions", "err", err)
							cancel()
							return
						}
					}
				}
			}
		}
	}()

	select {
	case <-ctx.Done():
		return -1, ctx.Err()
	case firstHeight := <-heightChan:
		return firstHeight, nil
	case <-time.After(maxSubscribeWait):
		return -1, fmt.Errorf("timeout: maxSubscribeWait")
	}
}

type counter struct {
	Blocks int
	Txs    int
}

func (s *Service) tryFetchAll(ctx context.Context, startHeight, endHeight, totalTxs int) error {
	var st, en int
	st = startHeight
	en = min(st+blockBatchSize, endHeight)

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

		err = s.db.CreateBlockList(ctx, resp.ToModel())
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
					n, err := s.syncTransactions(ctx,
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
				n, err := s.syncTransactions(ctx, st, en, 0, txsCount)
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
		if en > endHeight {
			break
		}
	}

	slog.Info("Synchronized blocks", "blocks", counter.Blocks, "txs", counter.Txs)
	return nil
}

func (s *Service) syncTransactions(ctx context.Context,
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

	// process transactions (handle events with SQS)
	go s.processEvents(ctx, resp.Transactions)

	// save to DB
	transactions, err := resp.ToModel()
	if err != nil {
		return n, err
	}
	err = s.db.CreateTransactionList(ctx, transactions)
	if err != nil {
		return n, err
	}
	return n, nil
}

func (s *Service) processEvents(ctx context.Context, transactions []indexer.Transaction) {
	for _, tx := range transactions {
		for _, event := range tx.Response.Events {
			if event.GnoEvent.Type != "Transfer" {
				continue
			}
			slog.Info("Transfer Event found!")
			evt, err := processEvent(event.GnoEvent)
			if err != nil {
				slog.Error("processEvent", "err", err)
				break
			}

			slog.Info("Valid Event found!, sending to SQS...")

			msgBody, err := json.Marshal(evt)
			if err != nil {
				slog.Error("json.Marshal", "err", err)
				break
			}
			err = s.sqs.SendMessage(ctx, sqs.Message{
				Body:      string(msgBody),
				CreatedAt: time.Now(),
			})
			if err != nil {
				slog.Error("SendMessage", "err", err)
				break
			}
		}
	}
}
