package sqs

import "time"

type Message struct {
	Body          string    `json:"body"`
	ReceiptHandle *string   `json:"receipt_handle"`
	CreatedAt     time.Time `json:"created_at"`
}
