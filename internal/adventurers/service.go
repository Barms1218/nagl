package adventurers

import (
	"context"
	"github.com/Barms1218/nagl/internal/database"
	"github.com/google/uuid"
)

type AdventurerService struct {
	store *database.Store
}

func NewAdventurerService(s *database.Store) *AdventurerService {
	return &AdventurerService{store: s}
}

func (s *AdventurerService) ListAdventurers(ctx context.Context, request GetMembersRequest) ([]GetAdventurersResponse, error) {
	params := database.GetGuildMembersParams{
		GuildID: database.UUIDToPgtype(request.GuildID),
		SortBy:  request.SortBy,
	}

	models, err := s.store.GetGuildMembers(ctx, params)
	if err != nil {
		return nil, err
	}

	response := make([]GetAdventurersResponse, 0, len(models))
	for _, model := range models {
		r := GetAdventurersResponse{
			ID:              model.ID,
			CurrentRank:     int(model.CurrentRank),
			CurrentActivity: string(model.CurrentActivity),
			Name:            model.Name.String,
			Role:            string(model.Role),
		}
		response = append(response, r)
	}

	return response, nil
}

func (s *AdventurerService) GetAdventurerDetails(ctx context.Context, id uuid.UUID) (DetailsResponse, error) {
	model, err := s.store.GetAdventurerDetails(ctx, id)
	if err != nil {
		return DetailsResponse{}, err
	}

	response := DetailsResponse{
		PartyID:         database.PgTypeToUUID(model.PartyID),
		Name:            model.Name.String,
		CurrentRank:     int(model.CurrentRank),
		Role:            string(model.Role),
		UpkeepCost:      int(model.UpkeepCost),
		CurrentActivity: string(model.CurrentActivity),
	}
	return response, nil
}

func (s *AdventurerService) GetAdventurersByGuild(ctx context.Context, guildID uuid.UUID) ([]GetAdventurersResponse, error) {
	models, err := s.store.GetAdventurersByGuild(ctx, database.UUIDToPgtype(guildID))
	if err != nil {
		return nil, err
	}

	response := make([]GetAdventurersResponse, 0, len(models))
	for _, model := range models {
		r := GetAdventurersResponse{
			ID:              model.ID,
			CurrentRank:     int(model.CurrentRank),
			CurrentActivity: string(model.CurrentActivity),
			Name:            model.Name.String,
			Role:            string(model.Role),
		}
		response = append(response, r)
	}
	return response, nil
}

func (s *AdventurerService) GetAdventurersWithStatus(ctx context.Context, request AdventurersWithStatusRequest) ([]GetAdventurersResponse, error) {
	params := database.GetAdventurersWithStatusParams{
		GuildID:         database.UUIDToPgtype(request.GuildID),
		CurrentActivity: database.ActivityEnum(request.CurrentActivity),
		SortBy:          request.SortBy,
	}
	models, err := s.store.GetAdventurersWithStatus(ctx, params)
	if err != nil {
		return nil, err
	}

	response := make([]GetAdventurersResponse, 0, len(models))
	for _, model := range models {
		r := GetAdventurersResponse{
			ID:              model.ID,
			CurrentRank:     int(model.CurrentRank),
			CurrentActivity: string(model.CurrentActivity),
			Name:            model.Name.String,
			Role:            string(model.Role),
		}
		response = append(response, r)
	}
	return response, nil
}
