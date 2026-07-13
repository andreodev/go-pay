package webhook

import (
	"context"

	"github.com/andreodev/go-pay/internal/risk"
)

type Service struct {
	repository     *Repository
	riskService    RiskService
	riskRepository RiskRepository
}

type RiskService interface {
	CalculateRisk(input risk.RiskRequest, ctx context.Context) (*risk.RiskResponse, error)
}

type RiskRepository interface {
	CreateRisk(ctx context.Context, r *risk.Risk) error
}

func NewService(repository *Repository, riskService RiskService, riskRepository RiskRepository) *Service {
	return &Service{
		repository:     repository,
		riskService:    riskService,
		riskRepository: riskRepository,
	}
}

func (s *Service) CreateWebhook(ctx context.Context, input WebhookRequest) (*WebhookRequest, error) {
	if err := validateWebhook(input); err != nil {
		return nil, err
	}

	exists, err := s.repository.ExistsByEventID(ctx, input.EventID)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, ErrEventAlreadyExists
	}

	if err := s.repository.CreateEvent(ctx, &input); err != nil {
		return nil, err
	}

	riskResponse, err := s.riskService.CalculateRisk(risk.RiskRequest{
		Amount:   input.Amount,
		IP:       input.IP,
		Email:    input.Customer.Email,
		DeviceID: input.DeviceID,
	}, ctx)
	if err != nil {
		return nil, err
	}

	riskEntity := risk.Risk{
		EventID:   input.EventID,
		PaymentID: input.PaymentID,
		Score:     riskResponse.Score,
		Level:     riskResponse.Level,
		Reasons:   riskResponse.Reasons,
	}

	if err := s.riskRepository.CreateRisk(ctx, &riskEntity); err != nil {
		return nil, err
	}

	return &input, nil
}
