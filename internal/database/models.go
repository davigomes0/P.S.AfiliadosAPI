package database

import (
	"time"
)

type Partner struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	ApiKey string `json:"api_key"`
}

type Conversion struct {
	ID            int       `json:"id"`
	TransactionID string    `json:"transaction_id"`
	PartnerID     int       `json:"partner_id"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateConversionRequest struct {
	TransactionID string  `json:"transaction_id"`
	Amount        float64 `json:"amount"`
}
