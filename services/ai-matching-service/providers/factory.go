package providers

import (
	"fmt"
	"os"
)

func NewLLMProvider() (LLMProvider, error) {
	provider := os.Getenv("LLM_PROVIDER")

	switch provider {
	case "ollama":
		return NewOllamaProvider(), nil
	case "gemini":
		return NewGeminiProvider(), nil
	case "groq":
		return NewGroqProvider(), nil
	default:
		return nil, fmt.Errorf("unknown LLM provider: %s", provider)
	}
}
