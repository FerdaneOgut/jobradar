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

type GroqProvider struct {
	apiKey string
	model  string
}

type groqRequest struct {
	Model    string        `json:"model"`
	Messages []groqMessage `json:"messages"`
}

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewGroqProvider() *GroqProvider {
	return &GroqProvider{
		apiKey: os.Getenv("GROQ_API_KEY"),
		model:  os.Getenv("GROQ_MODEL"),
	}
}

func (p *GroqProvider) Name() string {
	return "groq"
}

func (p *GroqProvider) Match(cvText string, jobTitle string, jobDescription string) (*MatchResponse, error) {
	prompt := fmt.Sprintf(`You are a job matching assistant.
Given the following CV and job description, return a match score from 0 to 100 and a brief reason.

CV:
%s

Job Title: %s
Job Description: %s

Respond ONLY in this exact JSON format, nothing else:
{"score": <number>, "reason": "<brief reason>"}`, cvText, jobTitle, jobDescription)

	reqBody, _ := json.Marshal(groqRequest{
		Model: p.model,
		Messages: []groqMessage{
			{Role: "user", Content: prompt},
		},
	})

	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("Groq raw response: %s", string(bodyBytes))

	var groqResp groqResponse
	if err := json.Unmarshal(bodyBytes, &groqResp); err != nil {
		return nil, err
	}

	if len(groqResp.Choices) == 0 {
		return nil, fmt.Errorf("no response from Groq")
	}

	text := strings.TrimSpace(groqResp.Choices[0].Message.Content)
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
