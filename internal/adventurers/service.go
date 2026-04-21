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
	var params database.ListRecruitableAdventurersParams
	if filters.Name != nil {
		params.Name = database.StringToPgtype(*filters.Name)

	}
	if filters.MinRank != nil {
		params.MinRank = database.IntToPgtype(*filters.MinRank)
	}
	if filters.MaxRank != nil {
		params.MaxRank = database.IntToPgtype(*filters.MaxRank)
	}
	if filters.Role != nil {
		params.Role = database.NullRoleEnum{RoleEnum: database.RoleEnum(*filters.Role)}
	}
	if filters.SortBy != nil {
		params.SortBy = *filters.SortBy
	} else {
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
		CurrentRank:     GetRankString(int(model.CurrentRank)),
		Role:            string(model.Role),
		UpkeepCost:      int(model.UpkeepCost),
		CurrentActivity: string(model.CurrentActivity),
	}
	return response, nil
}

func (s *AdventurerService) HireAdventurer(ctx context.Context, r SetAdventurerHiredRequest) error {
	params := database.SetAdventurerHiredParams{
		GuildID: database.UUIDToPgtype(r.GuildID),
		ID:      r.AdventurerID,
	}

	if err := s.store.SetAdventurerHired(ctx, params); err != nil {
		return err
	}

	return nil
}

func (s *AdventurerService) ListGuildMembers(ctx context.Context, guildID uuid.UUID, filters GuildMemberFilters) ([]ListMembersResponse, error) {
	var params database.ListGuildMembersParams

	params.GuildID = database.UUIDToPgtype(guildID)
	if filters.Name != nil {
		params.Name = database.StringToPgtype(*filters.Name)

	}
	if filters.MinRank != nil {
		params.MinRank = database.IntToPgtype(*filters.MinRank)
	}
	if filters.MaxRank != nil {
		params.MaxRank = database.IntToPgtype(*filters.MaxRank)
	}
	if filters.Role != nil {
		params.Role = database.NullRoleEnum{RoleEnum: database.RoleEnum(*filters.Role)}
	}
	if filters.SortBy != nil {
		params.SortBy = *filters.SortBy
	} else {
		params.SortBy = "name"
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

func (s *AdventurerService) GetUpkeepCost(ctx context.Context, adventurerID uuid.UUID) (int32, error) {
	return s.store.GetAdventurerUpkeepCost(ctx, adventurerID)
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
