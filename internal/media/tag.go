package media

type Tag struct {
	Value           string  `json:"value"`
	ConfidenceScore float32 `json:"confidence_score"`
	Provider        string  `json:"provider"`
}
