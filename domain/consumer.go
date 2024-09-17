package domain

import (
	"fbOnboarding/ent"
	"time"
)

type Consumer struct {
	CustomerID       int64     `json:"customer_id"`
	Username         string    `json:"username"`
	SubmissionStatus bool      `json:"submission_status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func ConsumerFromSchema(consumer *ent.Consumer) *Consumer {
	return &Consumer{
		CustomerID:       consumer.ID,
		Username:         consumer.Username,
		SubmissionStatus: consumer.SubmissionStatus,
		CreatedAt:        consumer.CreatedAt,
		UpdatedAt:        consumer.UpdatedAt,
	}
}
