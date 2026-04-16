package procedural

import (
	"context"
	"fmt"
	"github.com/anthropics/anthropic-sdk-go" // imported as anthropic
)

type ProceduralService struct {
	client *anthropic.Client
}

func NewProceduralService(c *anthropic.Client) *ProceduralService {
	return &ProceduralService{client: c}
}

func (s *ProceduralService) GenerateAdventurer(ctx context.Context) error {
	message, err := s.client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock("What is a quaternion?")),
		},
		Model: anthropic.ModelClaudeOpus4_6,
	})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v\n", message.Content)

	return nil
}
