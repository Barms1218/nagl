package procedural

import (
	"context"
	"fmt"
	"github.com/Barms1218/nagl/internal/database"

	"github.com/anthropics/anthropic-sdk-go" // imported as anthropic
)

type ProceduralService struct {
	client *anthropic.Client
	store  *database.Store
}

func NewProceduralService(c *anthropic.Client, s *database.Store) *ProceduralService {
	return &ProceduralService{
		client: c,
		store:  s,
	}
}

func (s *ProceduralService) GenerateAdventurer(ctx context.Context) (GeneratedContract, error) {
	// guild := s.store.GetGuild(ctx, guildID)
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

	return GeneratedContract{}, nil
}
