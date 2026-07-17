package webhook

import (
	"context"

	"github.com/andreodev/go-pay/internal/risk"
	"github.com/google/uuid"
)

type Service struct {
	repository     EventRepository
	riskService    RiskService
	riskRepository RiskRepository
}

type EventRepository interface {
	ExistsByEventID(ctx context.Context, eventID uuid.UUID) (bool, error)
	CreateEvent(ctx context.Context, input *WebhookRequest) error
}

type RiskService interface {
	CalculateRisk(input risk.RiskRequest, ctx context.Context) (*risk.RiskResponse, error)
}

type RiskRepository interface {
	CreateRisk(ctx context.Context, r *risk.Risk) error
}

func NewService(repository EventRepository, riskService RiskService, riskRepository RiskRepository) *Service {
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
		Amount:    input.Amount,
		IP:        input.IP,
		Email:     input.Customer.Email,
		Document:  input.Customer.Document,
		CardBin:   input.Card.Bin,
		CardLast4: input.Card.Last4,
		DeviceID:  input.DeviceID,
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
