package rules

const (
	DeviceIDRepeatedThreshold = 5
	DeviceIDRuleScore         = 30
)

func DeviceIDRule(count int, analysis *Analysis) {
	if count > DeviceIDRepeatedThreshold {
		analysis.Score += DeviceIDRuleScore
		analysis.Reasons = append(analysis.Reasons, "device_id_repeated")
	}
}
