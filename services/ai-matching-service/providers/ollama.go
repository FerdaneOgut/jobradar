package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type OllamaProvider struct {
	baseURL string
	model   string
}

type ollamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type ollamaResponse struct {
	Response string `json:"response"`
}

func NewOllamaProvider() *OllamaProvider {
	return &OllamaProvider{
		baseURL: os.Getenv("LLM_BASE_URL"),
		model:   os.Getenv("LLM_MODEL"),
	}
}

func (p *OllamaProvider) Name() string {
	return "ollama"
}

func (p *OllamaProvider) Match(cvText string, jobTitle string, jobDescription string) (*MatchResponse, error) {
	prompt := fmt.Sprintf(`You are a job matching assistant.
Given the following CV and job description, return a match score from 0 to 100 and a brief reason.

CV:
%s

Job Title: %s
Job Description: %s

Respond ONLY in this exact JSON format, nothing else:
{"score": <number>, "reason": "<brief reason>"}`, cvText, jobTitle, jobDescription)

	reqBody, _ := json.Marshal(ollamaRequest{
		Model:  p.model,
		Prompt: prompt,
		Stream: false,
	})

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(
		fmt.Sprintf("%s/api/generate", p.baseURL),
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ollamaResp ollamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, err
	}

	var matchResp MatchResponse
	if err := json.Unmarshal([]byte(ollamaResp.Response), &matchResp); err != nil {
		re := regexp.MustCompile(`"score"\s*:\s*(\d+)`)
		matches := re.FindStringSubmatch(ollamaResp.Response)
		if len(matches) > 1 {
			score, _ := strconv.Atoi(matches[1])
			matchResp.Score = score
			matchResp.Reason = strings.TrimSpace(ollamaResp.Response)
		}
	}

	return &matchResp, nil
}
