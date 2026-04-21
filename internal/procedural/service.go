package procedural

import (
	"context"
	"encoding/json"
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

func (s *ProceduralService) PromptForAdventurer(ctx context.Context) (*anthropic.Message, error) {
	systemPrompt, err := s.CreateAdventurerPrompt(ctx)
	if err != nil {
		return nil, err
	}
	message, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_6,
		MaxTokens: 512,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock("Generate a new adventurer available for recruitment."),
			),
		},
	})
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *ProceduralService) CreateAdventurerPrompt(ctx context.Context) (string, error) {
	existing, err := s.store.ListRecruitableAdventurers(ctx, database.ListRecruitableAdventurersParams{
		SortBy: "name",
	})
	if err != nil {
		return "", fmt.Errorf("Failed to fetch existing adventurers: %w", err)
	}
	type exclusionEntry struct {
		Name string `json:"name"`
		Role string `json:"role"`
		Rank int32  `json:"rank"`
	}
	exclusions := make([]exclusionEntry, 0, len(existing))
	for _, a := range existing {
		exclusion := exclusionEntry{
			Name: a.Name.String,
			Role: string(a.Role),
			Rank: a.CurrentRank,
		}
		exclusions = append(exclusions, exclusion)
	}

	exclusionJSON, err := json.Marshal(exclusions)
	if err != nil {
		return "", fmt.Errorf("Error serializing exclusion list: %w", err)
	}

	systemPrompt := fmt.Sprintf(`You are a fantasy adventurer record keeper. Generate an adventurer profile as JSON. Respond with ONLY valid JSON, no markdown, no explanation. Use this exact shape:
	{
		"name": "string",
		"role": "frontliner" | "spellcaster" | "healer" | "generalist",
		"current_rank": 1-5,
		"description": "string (2-3 sentences of flavor text)",
		"upkeep_cost": 10-100,
		"recruitment_cost": 50-(100 * current_rank)
	}
	The following adventurers already exist. Do not repeat their names, and avoid generating a duplicate combination of role and rank: %s`, string(exclusionJSON))

	return systemPrompt, nil
}

func (s *ProceduralService) GenerateAdventurer(ctx context.Context) (GeneratedAdventurer, error) {
	message, err := s.PromptForAdventurer(ctx)
	if err != nil {
		return GeneratedAdventurer{}, err
	}

	var raw string
	for _, block := range message.Content {
		if block.Type == "text" {
			raw = block.Text
			break
		}
	}
	if raw == "" {
		return GeneratedAdventurer{}, fmt.Errorf("Empty response from model")
	}

	var request AdventurerRequest
	if err := json.Unmarshal([]byte(raw), &request); err != nil {
		return GeneratedAdventurer{}, fmt.Errorf("JSON parse failed: %w - raw: %s", err, raw)
	}

	params := database.UpsertAdventurerParams{
		Name:        database.StringToPgtype(request.Name),
		CurrentRank: int32(request.Rank),
		Role:        database.RoleEnum(request.Role),
		Description: request.Bio,
	}

	adventurer := GeneratedAdventurer{
		Name:            request.Name,
		Role:            request.Role,
		Rank:            GetRankString(request.Rank),
		Bio:             request.Bio,
		UpkeepCost:      request.UpkeepCost,
		RecruitmentCost: request.RecruitmentCost,
	}

	_, err = s.store.UpsertAdventurer(ctx, params)
	if err != nil {
		return GeneratedAdventurer{}, fmt.Errorf("Error inserting new adventurer: %w", err)
	}

	return adventurer, nil
}

func (s *ProceduralService) CreateContractPrompt(ctx context.Context) (string, error) {
	existing, err := s.store.ListGuildContracts(ctx, database.ListGuildContractsParams{
		SortBy: "title",
	})
	if err != nil {
		return "", fmt.Errorf("Failed to fetch guild contracts: %w", err)
	}
	type exclusionEntry struct {
		Title        string `json:"title"`
		Difficulty   int32  `json:"difficulty" validate:"gte=1,lte=5"`
		RecPartySize int32  `json:"rec_party_size" validate:"gte=1,lte=5"`
	}

	exclusions := make([]exclusionEntry, 0, len(existing))
	for _, c := range existing {
		exclusion := exclusionEntry{
			Title:        c.Title.String,
			Difficulty:   c.Difficulty,
			RecPartySize: c.RecPartySize,
		}
		exclusions = append(exclusions, exclusion)
	}

	exclusionJSON, err := json.Marshal(exclusions)
	if err != nil {
		return "", fmt.Errorf("Error serializing exclusion list: %w", err)
	}

	systemPrompt := fmt.Sprintf(`You are a fantasy guild contract liason. Generate a contract as JSON. Response with ONLY valid JSON, no markdown, no explanation. Use this exact shape:
	{
		"title": "string",
		"description": string (3-4 sentences of flavor text related to the title.)",
		"difficulty": 1-5,
		"rec_party_size": 1-5,
		"reward": 300-500
	}
	The following contracts are already in use. Do not repeat their titles, and avoid generating a duplicate combination of difficulty and rec_party_size: %s`, string(exclusionJSON))

	return systemPrompt, nil
}

func (s *ProceduralService) PromptForContract(ctx context.Context) (*anthropic.Message, error) {
	systemPrompt, err := s.CreateAdventurerPrompt(ctx)
	if err != nil {
		return nil, err
	}
	message, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_6,
		MaxTokens: 512,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock("Generate a new contract."),
			),
		},
	})
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *ProceduralService) GenerateContract(ctx context.Context) (GeneratedContract, error) {
	message, err := s.PromptForContract(ctx)
	if err != nil {
		return GeneratedContract{}, err
	}

	var raw string
	for _, block := range message.Content {
		if block.Type == "text" {
			raw = block.Text
			break
		}
	}
	if raw == "" {
		return GeneratedContract{}, fmt.Errorf("Empty response model")
	}

	var request ContractRequest
	if err := json.Unmarshal([]byte(raw), &request); err != nil {
		return GeneratedContract{}, fmt.Errorf("JSON parse failed: %w - raw: %s", err, raw)
	}

	params := database.InsertContractParams{
		Title:        database.StringToPgtype(request.Title),
		Difficulty:   int32(request.Difficulty),
		RecPartySize: int32(request.RecPartySize),
		Description:  database.StringToPgtype(request.Description),
		Reward:       int32(request.Reward),
	}

	_, err = s.store.InsertContract(ctx, params)
	if err != nil {
		return GeneratedContract{}, fmt.Errorf("Failed to insert new contract: %w", err)
	}

	return GeneratedContract{
		Title:        request.Title,
		Description:  request.Description,
		Difficulty:   request.Description,
		RecPartySize: request.RecPartySize,
		Reward:       request.Reward,
	}, nil
}

func GetDifficultyString(difficulty int32) string {
	var diffString string
	switch difficulty := 1; difficulty {
	case 1:
		diffString = "Common"
	case 2:
		diffString = "Challening"
	case 3:
		diffString = "Dangerous"
	case 4:
		diffString = "Deadly"
	case 5:
		diffString = "Fatal"
	default:
		diffString = ""
	}
	return diffString
}

func GetRankString(rank int) string {
	var rankString string
	switch rank := 1; rank {
	case 1:
		rankString = "Iron"
	case 2:
		rankString = "Bronze"
	case 3:
		rankString = "Silver"
	case 4:
		rankString = "Gold"
	case 5:
		rankString = "Diamond"
	default:
		rankString = ""

	}
	return rankString
}
