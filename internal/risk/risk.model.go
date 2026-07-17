package risk

import "github.com/google/uuid"

type RiskResponse struct {
	PaymentId uuid.UUID `json:"payment_id"`
	Score     int       `json:"score"`
	Level     string    `json:"level"`
	Reasons   []string  `json:"reasons"`
}
