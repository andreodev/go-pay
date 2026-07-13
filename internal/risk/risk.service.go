package risk

import (
	"context"

	"github.com/andreodev/go-pay/internal/risk/rules"
	"github.com/google/uuid"
)

const (
	LowRiskMax    = 30
	MediumRiskMax = 70
)

type Analysis struct {
	Score   int
	Reasons []string
}

type RiskRequest struct {
	Amount int `json:"amount"`
}

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CalculateRisk(input RiskRequest) (*RiskResponse, error) {
	analysis := &rules.Analysis{}

	rules.AmountRule(input.Amount, analysis)

	level := calculateLevel(analysis.Score)

	return &RiskResponse{
		Score:   analysis.Score,
		Level:   level,
		Reasons: analysis.Reasons,
	}, nil
}

func calculateLevel(score int) string {
	if score <= LowRiskMax {
		return "LOW"
	}

	if score <= MediumRiskMax {
		return "MEDIUM"
	}

	return "HIGH"
}

func (s *Service) GetRiskPaymentID(paymentID uuid.UUID) (string, error) {
	risk, err := s.repository.GetRiskByPaymentID(context.Background(), paymentID)
	if err != nil {
		return "", err
	}

	return risk.Level, nil
}
