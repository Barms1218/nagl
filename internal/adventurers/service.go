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

func (s *AdventurerService) ListRecruitableAdventurers(ctx context.Context, filters SearchFilters) ([]ListAdventurersResponse, error) {
	params := database.ListRecruitableAdventurersParams{
		Name:    database.StringToPgtype(*filters.Name),
		MinRank: database.IntToPgtype(*filters.MinRank),
		MaxRank: database.IntToPgtype(*filters.MaxRank),
		Role:    database.NullRoleEnum{RoleEnum: database.RoleEnum(*filters.Role)},
		SortBy:  *filters.SortBy,
	}

	if params.SortBy == "" {
		params.SortBy = "name"
	}

	models, err := s.store.ListRecruitableAdventurers(ctx, params)
	if err != nil {
		return nil, err
	}

	response := make([]ListAdventurersResponse, 0, len(models))
	for _, model := range models {
		r := ListAdventurersResponse{
			ID:   model.ID,
			Rank: int32(model.CurrentRank),
			Name: model.Name.String,
			Role: string(model.Role),
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

func (s *AdventurerService) ListGuildMembers(ctx context.Context, guildID uuid.UUID, filters GuildMemberFilters) ([]ListMembersResponse, error) {
	params := database.ListGuildMembersParams{
		GuildID:         database.UUIDToPgtype(guildID),
		Name:            database.StringToPgtype(*filters.Name),
		MinRank:         database.IntToPgtype(*filters.MinRank),
		MaxRank:         database.IntToPgtype(*filters.MaxRank),
		Role:            database.NullRoleEnum{RoleEnum: database.RoleEnum(*filters.Role)},
		CurrentActivity: database.NullActivityEnum{ActivityEnum: database.ActivityEnum(*filters.Activity)},
		SortBy:          *filters.SortBy,
	}
	models, err := s.store.ListGuildMembers(ctx, params)
	if err != nil {
		return nil, err
	}

	response := make([]ListMembersResponse, 0, len(models))
	for _, model := range models {
		r := ListMembersResponse{
			ID:              model.ID,
			Rank:            int32(model.CurrentRank),
			CurrentActivity: string(model.CurrentActivity),
			Name:            model.Name.String,
			Role:            string(model.Role),
		}
		response = append(response, r)
	}
	return response, nil
}
