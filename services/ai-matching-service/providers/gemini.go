package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type GeminiProvider struct {
	apiKey string
	model  string
}

type geminiRequest struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func NewGeminiProvider() *GeminiProvider {
	return &GeminiProvider{
		apiKey: os.Getenv("GEMINI_API_KEY"),
		model:  os.Getenv("GEMINI_MODEL"),
	}
}

func (p *GeminiProvider) Name() string {
	return "gemini"
}

func (p *GeminiProvider) Match(cvText string, jobTitle string, jobDescription string) (*MatchResponse, error) {
	prompt := fmt.Sprintf(`You are a job matching assistant.
Given the following CV and job description, return a match score from 0 to 100 and a brief reason.

CV:
%s

Job Title: %s
Job Description: %s

Respond ONLY in this exact JSON format, nothing else:
{"score": <number>, "reason": "<brief reason>"}`, cvText, jobTitle, jobDescription)

	reqBody, _ := json.Marshal(geminiRequest{
		Contents: []geminiContent{
			{
				Parts: []geminiPart{
					{Text: prompt},
				},
			},
		},
	})

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		p.model,
		p.apiKey,
	)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read raw body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Gemini raw response: %s", string(bodyBytes))

	var geminiResp geminiResponse
	if err := json.Unmarshal(bodyBytes, &geminiResp); err != nil {
		return nil, err
	}

	if len(geminiResp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Gemini")
	}

	text := geminiResp.Candidates[0].Content.Parts[0].Text
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var matchResp MatchResponse
	if err := json.Unmarshal([]byte(text), &matchResp); err != nil {
		re := regexp.MustCompile(`"score"\s*:\s*(\d+)`)
		matches := re.FindStringSubmatch(text)
		if len(matches) > 1 {
			score, _ := strconv.Atoi(matches[1])
			matchResp.Score = score
			matchResp.Reason = text
		}
	}

	return &matchResp, nil
}
