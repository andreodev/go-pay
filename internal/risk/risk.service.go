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
	Amount    int       `json:"amount"`
	IP        string    `json:"ip"`
	Email     string    `json:"email"`
	Document  string    `json:"document"`
	CardBin   string    `json:"card_bin"`
	CardLast4 string    `json:"card_last4"`
	DeviceID  uuid.UUID `json:"device_id"`
}

type Service struct {
	repository RepositoryReader
}

type RepositoryReader interface {
	CountEventByIP(ctx context.Context, ip string) (int, error)
	CountDeclinedByEmail(ctx context.Context, email string) (int, error)
	CountOtherDocumentsByCard(ctx context.Context, bin string, last4 string, document string) (int, error)
	CountEventDeviceID(ctx context.Context, deviceID uuid.UUID) (int, error)
	GetRiskByPaymentID(ctx context.Context, paymentID uuid.UUID) (*Risk, error)
}

func NewService(repository RepositoryReader) *Service {
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

	cardDocumentCount, err := s.repository.CountOtherDocumentsByCard(ctx, input.CardBin, input.CardLast4, input.Document)
	if err != nil {
		return nil, err
	}

	rules.CardDocumentsRule(cardDocumentCount, analysis)

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
