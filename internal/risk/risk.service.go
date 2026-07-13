package risk

import (
	"context"

	"github.com/google/uuid"
)

const (
	HighAmountThreshold = 100_000

	AmountRuleScore = 20

	LowRiskMax    = 30
	MediumRiskMax = 70
)

type RiskRequest struct {
	Amount int `json:"amount"`
}

type Service struct {
	repository     *Repository
	riskRepository RiskRepository
}

type RiskRepository interface {
	GetRiskByPaymentID(ctx context.Context, paymentID uuid.UUID) (*Risk, error)
}

func NewService(repository *Repository, riskRepository RiskRepository) *Service {
	return &Service{
		repository:     repository,
		riskRepository: riskRepository,
	}
}

func (s *Service) CalculateRisk(input RiskRequest) (*RiskResponse, error) {
	score := 0
	reasons := []string{}

	if input.Amount > HighAmountThreshold {
		score += AmountRuleScore
		reasons = append(reasons, "amount_above_1000")
	}

	level := calculateLevel(score)

	return &RiskResponse{
		Score:   score,
		Level:   level,
		Reasons: reasons,
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
	risk, err := s.riskRepository.GetRiskByPaymentID(context.Background(), paymentID)
	if err != nil {
		return "", err
	}

	return risk.Level, nil
}
