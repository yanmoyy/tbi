package synchronizer

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/yanmoyy/tbi/internal/database"
	"github.com/yanmoyy/tbi/internal/indexer"
)

const (
	maxBatchSize = 10
	maxRetries   = 3
	retryDelay   = time.Second * 5
	timeout      = time.Second * 5
)

type Service struct {
	indexer   *indexer.Client
	db        *database.Client
	sqsClient *sqs.Client
}

func New(client *indexer.Client, db *database.Client, sqsClient *sqs.Client) *Service {
	return &Service{
		indexer:   client,
		db:        db,
		sqsClient: sqsClient,
	}
}
