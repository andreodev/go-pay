package payments

import "context"

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) ListPaymentEvents(
	ctx context.Context,
	page int,
	limit int,
) (*PaymentEventListResponse, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	events, total, err := s.repository.ListPaymentsEvents(
		ctx,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	totalPages := 0

	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	return &PaymentEventListResponse{
		Data:        events,
		Page:        page,
		Limit:       limit,
		Total:       total,
		Total_Pages: totalPages,
	}, nil
}
