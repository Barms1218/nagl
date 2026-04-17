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

func (s *AdventurerService) ListAdventurers(ctx context.Context, request GetMembersRequest) ([]GetMembersResponse, error) {
	params := database.GetGuildMembersParams{
		GuildID: request.GuildID,
		SortBy:  request.SortBy,
	}

	models, err := s.store.GetGuildMembers(ctx, params)
	if err != nil {
		return nil, err
	}

	response := make([]GetMembersResponse, 0, len(models))
	for _, model := range models {
		r := GetMembersResponse{
			ID:              database.PgTypeToUUID(model.ID),
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
	model, err := s.store.GetAdventurerDetails(ctx, database.UUIDToPgtype(id))
	if err != nil {
		return DetailsResponse{}, err
	}

	response := DetailsResponse{
		PartyID:         model.PartyID,
		Name:            model.Name.String,
		CurrentRank:     int(model.CurrentRank),
		Role:            string(model.Role),
		UpkeepCost:      int(model.UpkeepCost),
		CurrentActivity: string(model.CurrentActivity),
	}
	return response, nil
}
