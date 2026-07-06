package clients

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

var ErrNoAPIKey = errors.New("OPENAI_API_KEY not set")

const systemPrompt = "You are a personal finance advisor. Respond ONLY with valid JSON matching the exact schema provided. No markdown, no code fences, no extra text."

func Configured() bool {
	return os.Getenv("OPENAI_API_KEY") != ""
}

func Chat(ctx context.Context, prompt string) (string, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return "", ErrNoAPIKey
	}

	client := openai.NewClient(key)
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	resp, err := client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Temperature: 0.7,
	})
	if err != nil {
		return "", fmt.Errorf("AI request failed: %w", err)
	}
	return resp.Choices[0].Message.Content, nil
}
