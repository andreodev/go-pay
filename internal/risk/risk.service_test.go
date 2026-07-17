package risk

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

type fakeRiskRepository struct {
	ipCount            int
	emailCount         int
	cardDocumentsCount int
	deviceIDCount      int
}

func (f fakeRiskRepository) CountEventByIP(context.Context, string) (int, error) {
	return f.ipCount, nil
}

func (f fakeRiskRepository) CountDeclinedByEmail(context.Context, string) (int, error) {
	return f.emailCount, nil
}

func (f fakeRiskRepository) CountOtherDocumentsByCard(context.Context, string, string, string) (int, error) {
	return f.cardDocumentsCount, nil
}

func (f fakeRiskRepository) CountEventDeviceID(context.Context, uuid.UUID) (int, error) {
	return f.deviceIDCount, nil
}

func (f fakeRiskRepository) GetRiskByPaymentID(context.Context, uuid.UUID) (*Risk, error) {
	return nil, nil
}

func TestCalculateLevel(t *testing.T) {
	tests := []struct {
		score int
		want  string
	}{
		{0, "LOW"},
		{LowRiskMax, "LOW"},
		{LowRiskMax + 1, "MEDIUM"},
		{MediumRiskMax, "MEDIUM"},
		{MediumRiskMax + 1, "HIGH"},
	}

	for _, tt := range tests {
		if got := calculateLevel(tt.score); got != tt.want {
			t.Fatalf("calculateLevel(%d) = %s, want %s", tt.score, got, tt.want)
		}
	}
}

func TestCalculateRisk(t *testing.T) {
	deviceID := uuid.New()
	tests := []struct {
		name       string
		request    RiskRequest
		repository fakeRiskRepository
		wantScore  int
		wantLevel  string
		wantReason []string
	}{
		{
			name: "low risk with no triggered rules",
			request: RiskRequest{
				Amount:    1000,
				IP:        "10.0.0.1",
				Email:     "user@example.com",
				Document:  "11111111111",
				CardBin:   "411111",
				CardLast4: "1111",
				DeviceID:  deviceID,
			},
			wantScore: 0,
			wantLevel: "LOW",
		},
		{
			name: "medium risk from amount and email",
			request: RiskRequest{
				Amount:    100_001,
				IP:        "10.0.0.1",
				Email:     "user@example.com",
				Document:  "11111111111",
				CardBin:   "411111",
				CardLast4: "1111",
				DeviceID:  deviceID,
			},
			repository: fakeRiskRepository{emailCount: 3},
			wantScore:  45,
			wantLevel:  "MEDIUM",
			wantReason: []string{"High amount transaction", "email_many_declines"},
		},
		{
			name: "high risk when all rules trigger",
			request: RiskRequest{
				Amount:    100_001,
				IP:        "10.0.0.1",
				Email:     "user@example.com",
				Document:  "11111111111",
				CardBin:   "411111",
				CardLast4: "1111",
				DeviceID:  deviceID,
			},
			repository: fakeRiskRepository{
				ipCount:            6,
				emailCount:         3,
				cardDocumentsCount: 1,
				deviceIDCount:      6,
			},
			wantScore: 145,
			wantLevel: "HIGH",
			wantReason: []string{
				"High amount transaction",
				"ip_repeated",
				"email_many_declines",
				"card_many_documents",
				"device_id_repeated",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(tt.repository)

			got, err := service.CalculateRisk(tt.request, context.Background())
			if err != nil {
				t.Fatal(err)
			}

			if got.Score != tt.wantScore {
				t.Fatalf("score = %d, want %d", got.Score, tt.wantScore)
			}
			if got.Level != tt.wantLevel {
				t.Fatalf("level = %s, want %s", got.Level, tt.wantLevel)
			}
			if !reflect.DeepEqual(got.Reasons, tt.wantReason) {
				t.Fatalf("reasons = %v, want %v", got.Reasons, tt.wantReason)
			}
		})
	}
}
