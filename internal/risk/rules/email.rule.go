package rules

const (
	EmailDeclinedThreshold = 3
	EmailRuleScore         = 25
)

func EmailRule(count int, analysis *Analysis) {
	if count >= EmailDeclinedThreshold {
		analysis.Score += EmailRuleScore
		analysis.Reasons = append(analysis.Reasons, "email_many_declines")
	}
}
