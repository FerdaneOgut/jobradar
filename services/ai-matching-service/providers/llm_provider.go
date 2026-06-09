package providers

type MatchResponse struct {
	Score  int    `json:"score"`
	Reason string `json:"reason"`
}

type LLMProvider interface {
	Match(cvText string, jobTitle string, jobDescription string) (*MatchResponse, error)
	Name() string
}
