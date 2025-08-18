package sqs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/yanmoyy/tbi/internal/config"
)

type Client struct {
	sqs      *sqs.Client
	queueUrl string
}

func NewClient(cfg config.SQS) *Client {
	return &Client{
		sqs: sqs.NewFromConfig(aws.Config{
			BaseEndpoint: aws.String(cfg.Endpoint),
		}),
		queueUrl: cfg.QueueURL,
	}
}

func (c *Client) SendMessage(ctx context.Context, msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}
	input := &sqs.SendMessageInput{
		QueueUrl:    aws.String(c.queueUrl),
		MessageBody: aws.String(string(body)),
	}
	_, err = c.sqs.SendMessage(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}
	return nil
}

func (c *Client) GetMessages(ctx context.Context, maxMessages int32, waitTime int32) ([]Message, error) {
	var messages []Message
	result, err := c.sqs.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueUrl),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     waitTime,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to receive message: %w", err)
	} else {
		messages = make([]Message, len(result.Messages))
		for i, m := range result.Messages {
			var msg Message
			err := json.Unmarshal([]byte(*m.Body), &msg)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal message: %w", err)
			}
			msg.ReceiptHandle = m.ReceiptHandle
			messages[i] = msg
		}
	}
	return messages, nil
}

func (c *Client) DeleteMessages(ctx context.Context, messages []Message) error {
	entries := make([]types.DeleteMessageBatchRequestEntry, len(messages))
	for msgIndex := range messages {
		entries[msgIndex].Id = aws.String(fmt.Sprintf("%v", msgIndex))
		entries[msgIndex].ReceiptHandle = messages[msgIndex].ReceiptHandle
	}
	_, err := c.sqs.DeleteMessageBatch(ctx, &sqs.DeleteMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(c.queueUrl),
	})
	if err != nil {
		log.Printf("Couldn't delete messages from queue %v. Here's why: %v\n", c.queueUrl, err)
	}
	return nil
}

func (c *Client) PurgeQueue(ctx context.Context) error {
	_, err := c.sqs.PurgeQueue(ctx, &sqs.PurgeQueueInput{
		QueueUrl: aws.String(c.queueUrl)})
	if err != nil {
		log.Printf("Couldn't purge queue %v. Here's why: %v\n", c.queueUrl, err)
	}
	return nil
}
