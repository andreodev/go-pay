package rules

import "testing"

func TestAmountRule(t *testing.T) {
	tests := []struct {
		name       string
		amount     int
		wantScore  int
		wantReason bool
	}{
		{"at threshold", HighAmountThreshold, 0, false},
		{"above threshold", HighAmountThreshold + 1, AmountRuleScore, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &Analysis{}

			AmountRule(tt.amount, analysis)

			if analysis.Score != tt.wantScore {
				t.Fatalf("score = %d, want %d", analysis.Score, tt.wantScore)
			}
			if (len(analysis.Reasons) > 0) != tt.wantReason {
				t.Fatalf("reasons = %v, want reason: %t", analysis.Reasons, tt.wantReason)
			}
		})
	}
}

func TestIPRule(t *testing.T) {
	tests := []struct {
		name       string
		count      int
		wantScore  int
		wantReason bool
	}{
		{"at threshold", IPRepeatedThreshold, 0, false},
		{"above threshold", IPRepeatedThreshold + 1, IPRuleScore, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &Analysis{}

			IPRule(tt.count, analysis)

			if analysis.Score != tt.wantScore {
				t.Fatalf("score = %d, want %d", analysis.Score, tt.wantScore)
			}
			if (len(analysis.Reasons) > 0) != tt.wantReason {
				t.Fatalf("reasons = %v, want reason: %t", analysis.Reasons, tt.wantReason)
			}
		})
	}
}

func TestEmailRule(t *testing.T) {
	tests := []struct {
		name       string
		count      int
		wantScore  int
		wantReason bool
	}{
		{"below threshold", EmailDeclinedThreshold - 1, 0, false},
		{"at threshold", EmailDeclinedThreshold, EmailRuleScore, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &Analysis{}

			EmailRule(tt.count, analysis)

			if analysis.Score != tt.wantScore {
				t.Fatalf("score = %d, want %d", analysis.Score, tt.wantScore)
			}
			if (len(analysis.Reasons) > 0) != tt.wantReason {
				t.Fatalf("reasons = %v, want reason: %t", analysis.Reasons, tt.wantReason)
			}
		})
	}
}

func TestDeviceIDRule(t *testing.T) {
	tests := []struct {
		name       string
		count      int
		wantScore  int
		wantReason bool
	}{
		{"at threshold", DeviceIDRepeatedThreshold, 0, false},
		{"above threshold", DeviceIDRepeatedThreshold + 1, DeviceIDRuleScore, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &Analysis{}

			DeviceIDRule(tt.count, analysis)

			if analysis.Score != tt.wantScore {
				t.Fatalf("score = %d, want %d", analysis.Score, tt.wantScore)
			}
			if (len(analysis.Reasons) > 0) != tt.wantReason {
				t.Fatalf("reasons = %v, want reason: %t", analysis.Reasons, tt.wantReason)
			}
		})
	}
}

func TestCardDocumentsRule(t *testing.T) {
	tests := []struct {
		name       string
		count      int
		wantScore  int
		wantReason bool
	}{
		{"same document only", 0, 0, false},
		{"other document found", 1, CardDocumentsRuleScore, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis := &Analysis{}

			CardDocumentsRule(tt.count, analysis)

			if analysis.Score != tt.wantScore {
				t.Fatalf("score = %d, want %d", analysis.Score, tt.wantScore)
			}
			if (len(analysis.Reasons) > 0) != tt.wantReason {
				t.Fatalf("reasons = %v, want reason: %t", analysis.Reasons, tt.wantReason)
			}
		})
	}
}
