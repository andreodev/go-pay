package webhook

import (
	"context"
	"errors"
	"testing"

	"github.com/andreodev/go-pay/internal/risk"
	"github.com/google/uuid"
)

type fakeEventRepository struct {
	exists       bool
	createCount  int
	createdEvent *WebhookRequest
}

func (f *fakeEventRepository) ExistsByEventID(context.Context, uuid.UUID) (bool, error) {
	return f.exists, nil
}

func (f *fakeEventRepository) CreateEvent(_ context.Context, input *WebhookRequest) error {
	f.createCount++
	f.createdEvent = input
	return nil
}

type fakeRiskService struct {
	calculateCount int
	request        risk.RiskRequest
}

func (f *fakeRiskService) CalculateRisk(input risk.RiskRequest, _ context.Context) (*risk.RiskResponse, error) {
	f.calculateCount++
	f.request = input
	return &risk.RiskResponse{
		Score:   40,
		Level:   "MEDIUM",
		Reasons: []string{"card_many_documents"},
	}, nil
}

type fakeRiskWriter struct {
	createCount int
	createdRisk *risk.Risk
}

func (f *fakeRiskWriter) CreateRisk(_ context.Context, r *risk.Risk) error {
	f.createCount++
	f.createdRisk = r
	return nil
}

func validWebhookRequest() WebhookRequest {
	return WebhookRequest{
		EventID:   uuid.New(),
		Type:      "payment.created",
		PaymentID: uuid.New(),
		Amount:    12000,
		Status:    "approved",
		Currency:  "BRL",
		Customer: Customer{
			Email:    "customer@example.com",
			Document: "11111111111",
		},
		Card: Card{
			Bin:   "411111",
			Last4: "1111",
		},
		IP:       "10.0.0.1",
		DeviceID: uuid.New(),
	}
}

func TestCreateWebhookRejectsDuplicateEvent(t *testing.T) {
	events := &fakeEventRepository{exists: true}
	risks := &fakeRiskWriter{}
	riskService := &fakeRiskService{}
	service := NewService(events, riskService, risks)

	got, err := service.CreateWebhook(context.Background(), validWebhookRequest())

	if !errors.Is(err, ErrEventAlreadyExists) {
		t.Fatalf("err = %v, want %v", err, ErrEventAlreadyExists)
	}
	if got != nil {
		t.Fatalf("webhook = %v, want nil", got)
	}
	if events.createCount != 0 {
		t.Fatalf("CreateEvent called %d times, want 0", events.createCount)
	}
	if riskService.calculateCount != 0 {
		t.Fatalf("CalculateRisk called %d times, want 0", riskService.calculateCount)
	}
	if risks.createCount != 0 {
		t.Fatalf("CreateRisk called %d times, want 0", risks.createCount)
	}
}

func TestCreateWebhookPersistsEventAndRisk(t *testing.T) {
	events := &fakeEventRepository{}
	risks := &fakeRiskWriter{}
	riskService := &fakeRiskService{}
	service := NewService(events, riskService, risks)
	input := validWebhookRequest()

	got, err := service.CreateWebhook(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}

	if got == nil || got.EventID != input.EventID {
		t.Fatalf("webhook = %v, want event %s", got, input.EventID)
	}
	if events.createCount != 1 {
		t.Fatalf("CreateEvent called %d times, want 1", events.createCount)
	}
	if riskService.request.Document != input.Customer.Document ||
		riskService.request.CardBin != input.Card.Bin ||
		riskService.request.CardLast4 != input.Card.Last4 {
		t.Fatalf("risk request = %+v, want customer document and card data", riskService.request)
	}
	if risks.createdRisk == nil {
		t.Fatal("risk was not created")
	}
	if risks.createdRisk.EventID != input.EventID || risks.createdRisk.PaymentID != input.PaymentID {
		t.Fatalf("risk = %+v, want event/payment IDs from webhook", risks.createdRisk)
	}
	if risks.createdRisk.Score != 40 || risks.createdRisk.Level != "MEDIUM" {
		t.Fatalf("risk = %+v, want risk service result", risks.createdRisk)
	}
}
