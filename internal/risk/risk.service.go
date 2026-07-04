package risk

const (
	HighAmountThreshold = 100_000

	AmountRuleScore = 20

	LowRiskMax    = 30
	MediumRiskMax = 70
)

type RiskRequest struct {
	Amount int `json:"amount"`
}

type Service struct{}

func NewService() *Service {
	return &Service{}
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
