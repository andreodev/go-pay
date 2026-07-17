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
	Amount   int       `json:"amount"`
	IP       string    `json:"ip"`
	Email    string    `json:"email"`
	DeviceID uuid.UUID `json:"device_id"`
}

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CalculateRisk(input RiskRequest, ctx context.Context) (*RiskResponse, error) {
	analysis := &rules.Analysis{}

	rules.AmountRule(input.Amount, analysis)

	count, err := s.repository.CountEventByIP(ctx, input.IP)
	if err != nil {
		return nil, err
	}

	rules.IPRule(count, analysis)

	emailCount, err := s.repository.CountDeclinedByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	rules.EmailRule(emailCount, analysis)

	deviceIDCount, err := s.repository.CountEventDeviceID(ctx, input.DeviceID)

	if err != nil {
		return nil, err
	}

	rules.DeviceIDRule(deviceIDCount, analysis)

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

func (s *Service) GetRiskPaymentID(paymentID uuid.UUID) (*RiskResponse, error) {
	risk, err := s.repository.GetRiskByPaymentID(context.Background(), paymentID)
	if err != nil {
		return nil, err
	}

	return &RiskResponse{
		PaymentId: risk.PaymentID,
		Score:     risk.Score,
		Level:     risk.Level,
		Reasons:   risk.Reasons,
	}, nil
}
