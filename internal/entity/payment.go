package entity

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	CorrelationId uuid.UUID `json:"correlationId"`
	Amount        float64   `json:"amount"`
	RequestedAt   time.Time `json:"requestedAt"`
	Processor     string    `json:"processor"` // "default" ou "fallback"
}
