package webhook

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
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

	return &input, nil
}
