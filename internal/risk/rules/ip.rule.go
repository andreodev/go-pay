package rules

const (
	IPRepeatedThreshold = 5
	IPRuleScore         = 30
)

func IPRule(count int, analysis *Analysis) {
	if count > IPRepeatedThreshold {
		analysis.Score += IPRuleScore
		analysis.Reasons = append(analysis.Reasons, "ip_repeated")
	}
}
