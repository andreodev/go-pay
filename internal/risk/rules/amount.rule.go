package rules

const (
	HighAmountThreshold = 100_000
	AmountRuleScore     = 20
)

type Analysis struct {
	Score   int
	Reasons []string
}

func AmountRule(amount int, analysis *Analysis) {
	if amount > HighAmountThreshold {
		analysis.Score += AmountRuleScore
		analysis.Reasons = append(analysis.Reasons, "High amount transaction")
	}
}
