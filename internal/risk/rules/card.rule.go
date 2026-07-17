package rules

const CardDocumentsRuleScore = 40

func CardDocumentsRule(otherDocumentsCount int, analysis *Analysis) {
	if otherDocumentsCount > 0 {
		analysis.Score += CardDocumentsRuleScore
		analysis.Reasons = append(analysis.Reasons, "card_many_documents")
	}
}
