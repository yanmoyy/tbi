package sqs

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yanmoyy/tbi/internal/config"
	"github.com/yanmoyy/tbi/internal/test"
)

var queueName = flag.String("queue", "test-queue", "queue name")

func getTestClient(t *testing.T) *Client {
	t.Helper()
	test.CheckSQSFlag(t)
	c := NewClient(config.SQS{Endpoint: "http://localhost:4566", QueueURL: "http://localhost:4566/000000000000/" + *queueName})
	return c
}

func TestSQSMessage(t *testing.T) {
	t.Run("send and receive", func(t *testing.T) {
		_, cancel := context.WithTimeout(t.Context(), 1*time.Second)
		defer cancel()

		c := getTestClient(t)
		defer c.PurgeQueue(t.Context())

		ctx := context.Background()
		//
		now := time.Now()
		msg := Message{
			Body:      "TestMessage",
			CreatedAt: now,
		}

		err := c.SendMessage(ctx, msg)
		require.NoError(t, err)

		messages, err := c.GetMessages(ctx, 1, 1)
		require.NoError(t, err)

		t.Logf("messages: %+v", messages)
		require.Len(t, messages, 1)
		require.Equal(t, msg.Body, messages[0].Body)
		require.Equal(t, now.Format(time.RFC3339), messages[0].CreatedAt.Format(time.RFC3339))

		err = c.DeleteMessages(ctx, []Message{messages[0]})
		require.NoError(t, err)
	})
}
