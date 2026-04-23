package procedural

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/invopop/jsonschema"

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

func GenerateSchema(v any) map[string]any {
	r := jsonschema.Reflector{AllowAdditionalProperties: false, DoNotReference: true}
	s := r.Reflect(v)
	b, _ := json.Marshal(s)
	var m map[string]any
	json.Unmarshal(b, &m)
	return m
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
		OutputConfig: anthropic.OutputConfigParam{
			Format: anthropic.JSONOutputFormatParam{
				Schema: GenerateSchema(GeneratedAdventurer{}),
			},
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

	systemPrompt := fmt.Sprintf(`You are a fantasy adventurer record keeper. Generate an adventurer profile as JSON. 

	The following adventurers already exist. Do not repeat their names, and avoid generating a duplicate combination of role and rank: %s`, string(exclusionJSON))

	return systemPrompt, nil
}

func (s *ProceduralService) GenerateAdventurer(ctx context.Context) error {
	message, err := s.PromptForAdventurer(ctx)
	if err != nil {
		return err
	}

	var adventurer GeneratedAdventurer
	for _, block := range message.Content {
		switch variant := block.AsAny().(type) {
		case anthropic.TextBlock:
			if err := json.Unmarshal([]byte(variant.Text), &adventurer); err != nil {
				return fmt.Errorf("JSON parse failed: %w - raw: %s", err, variant.Text)
			}
		}
	}
	params := database.UpsertAdventurerParams{
		Name:        database.StringToPgtype(adventurer.Name),
		CurrentRank: adventurer.Rank,
		Role:        database.RoleEnum(adventurer.Role),
		Description: adventurer.Bio,
	}

	_, err = s.store.UpsertAdventurer(ctx, params)
	if err != nil {
		return fmt.Errorf("Error inserting new adventurer: %w", err)
	}

	return nil
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

	systemPrompt := fmt.Sprintf(`You are a fantasy guild contract liason. Generate a contract as JSON.

	The following contracts are already in use. Do not repeat their titles, and avoid generating a duplicate combination of difficulty and rec_party_size: %s`, string(exclusionJSON))

	return systemPrompt, nil
}

func (s *ProceduralService) PromptForContract(ctx context.Context) (*anthropic.Message, error) {
	systemPrompt, err := s.CreateContractPrompt(ctx)
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
		OutputConfig: anthropic.OutputConfigParam{
			Format: anthropic.JSONOutputFormatParam{
				Schema: GenerateSchema(GeneratedContract{}),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *ProceduralService) GenerateContract(ctx context.Context) error {
	message, err := s.PromptForContract(ctx)
	if err != nil {
		return err
	}

	var contract GeneratedContract
	for _, block := range message.Content {
		switch variant := block.AsAny().(type) {
		case anthropic.TextBlock:
			if err := json.Unmarshal([]byte(variant.Text), &contract); err != nil {
				return fmt.Errorf("JSON parse failed: %w - raw: %s", err, variant.Text)
			}
		}
	}

	params := database.InsertContractParams{
		Title:           database.StringToPgtype(contract.Title),
		Difficulty:      contract.Difficulty,
		RecPartySize:    int32(contract.RecPartySize),
		Description:     database.StringToPgtype(contract.Description),
		Reward:          int32(contract.Reward),
		DurationMinutes: int32(contract.Duration),
	}

	_, err = s.store.InsertContract(ctx, params)
	if err != nil {
		return fmt.Errorf("Failed to insert new contract: %w", err)
	}

	return nil
}

func (s *ProceduralService) GenerateParty(ctx context.Context, r GeneratePartyRequest) (GeneratedParty, error) {

	partyJSON, err := json.Marshal(r)
	if err != nil {
		return GeneratedParty{}, fmt.Errorf("Failed to serialize party info: %w", err)
	}
	systemPrompt := fmt.Sprintf(`You are a fantasy guild party manager. Generate a name for a party of fantasy adventurers. The party details are as follows: %s`, string(partyJSON))

	message, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_6,
		MaxTokens: 512,
		System: []anthropic.TextBlockParam{
			{Text: systemPrompt},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock("Generate a new party name."),
			),
		},
		OutputConfig: anthropic.OutputConfigParam{
			Format: anthropic.JSONOutputFormatParam{
				Schema: GenerateSchema(PartyName{}),
			},
		},
	})
	if err != nil {
		return GeneratedParty{}, err
	}

	var party PartyName
	for _, block := range message.Content {
		switch variant := block.AsAny().(type) {
		case anthropic.TextBlock:
			if err := json.Unmarshal([]byte(variant.Text), &party); err != nil {
				return GeneratedParty{}, fmt.Errorf("JSON parse failed: %w - raw: %s", err, variant.Text)
			}
		}
	}

	params := database.CreatePartyParams{
		GuildID: r.GuildID,
		Name:    party.PartyName,
	}
	inserted, err := s.store.CreateParty(ctx, params)

	return GeneratedParty{
		GuildID:      party.GuildID,
		GuildName:    r.GuildName,
		PartyName:    party.PartyName,
		PartyStatus:  string(inserted.PartyStatus),
		MaxPartySize: inserted.MaximumPartySize,
	}, nil
}

func GetDifficultyString(difficulty int32) string {
	var diffString string
	switch difficulty {
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
	switch rank {
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
