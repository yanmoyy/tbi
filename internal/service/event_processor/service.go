package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/yanmoyy/tbi/internal/database"
	"github.com/yanmoyy/tbi/internal/models"
	"github.com/yanmoyy/tbi/internal/sqs"
)

type Service struct {
	db  *database.Client
	sqs *sqs.Client
}

func New(db *database.Client, sqs *sqs.Client) *Service {
	return &Service{
		db:  db,
		sqs: sqs,
	}
}

func (s *Service) Run(ctx context.Context) error {
	slog.Info("start listening")
	return s.startListening(ctx)
}

func (s *Service) startListening(ctx context.Context) error {
	const pollInterval = 3 * time.Second
	for {
		select {
		case <-ctx.Done():
			slog.Info("Listening stopped")
			return ctx.Err()
		default:
			msgs, err := s.sqs.GetMessages(ctx, 1, 1)
			if err != nil {
				slog.Error("GetMessages", "error", err)
				goto interval
			}
			if err := s.processMessages(ctx, msgs); err != nil {
				slog.Error("Failed to process messages", "error", err)
				goto interval
			}
			if err := s.sqs.DeleteMessages(ctx, msgs); err != nil {
				slog.Error("sqs delete messages", "error", err)
				goto interval
			}
		interval:
			slog.Info("Waiting for " + pollInterval.String())
			time.Sleep(pollInterval)
		}
	}
}

func (s *Service) processMessages(ctx context.Context, msgs []sqs.Message) error {
	for _, msg := range msgs {
		var evt models.TransferEvent
		if err := json.Unmarshal([]byte(msg.Body), &evt); err != nil {
			slog.Error("unmarshal event", "error", err)
			continue
		}
		if err := s.handleEvent(ctx, evt); err != nil {
			slog.Error("handle event", "error", err)
			continue
		}
	}
	return nil
}

// handleEvent processes a single event and updates the database
func (s *Service) handleEvent(ctx context.Context, evt models.TransferEvent) error {
	slog.Info("Handling Event...", "event", evt)
	switch evt.Func {
	case models.Mint:
		err := s.db.UpdateTokenBalance(ctx, evt.To, evt.TokenPath, evt.Value, true)
		if err != nil {
			return fmt.Errorf("UpdateTokenBalance(Mint): %w", err)
		}
	case models.Burn:
		err := s.db.UpdateTokenBalance(ctx, evt.From, evt.TokenPath, evt.Value, false)
		if err != nil {
			return fmt.Errorf("UpdateTokenBalance(Burn): %w", err)
		}
	case models.Transfer:
		err := s.db.TransferTokenBalance(ctx, evt.From, evt.To, evt.TokenPath, evt.Value)
		if err != nil {
			return fmt.Errorf("TransferTokenBalance: %w", err)
		}
		err = s.db.CreateTokenTransfer(ctx, models.TokenTransfer{
			FromAddress: evt.From,
			ToAddress:   evt.To,
			TokenPath:   evt.TokenPath,
			Amount:      evt.Value,
		})
		if err != nil {
			return fmt.Errorf("CreateTokenTransfer: %w", err)
		}
	default:
		return fmt.Errorf("unsupported event function: %s", evt.Func)
	}
	return nil
}
