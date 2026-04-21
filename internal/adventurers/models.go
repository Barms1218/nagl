package adventurers

import (
	"github.com/google/uuid"
)

type SearchFilters struct {
	Name    *string `json:"name"`
	MinRank *int32  `json:"min_rank" validate:"gte=1,lte=5"`
	MaxRank *int32  `json:"max_rank" validate:"gte=1,lte=5"`
	Role    *string `json:"role" validate:"oneof=frontliner spellcaster healer generalist"`
	SortBy  *string `json:"sort_by" validate:"oneof=name current_rank role"`
}

type GuildMemberFilters struct {
	Name     *string `json:"name"`
	MinRank  *int32  `json:"min_rank" validate:"gte=1,lte=5"`
	MaxRank  *int32  `json:"max_rank" validate:"gte=1,lte=5"`
	Role     *string `json:"role" validate:"oneof=frontliner spellcaster healer generalist"`
	Activity *string `json:"current_activity" validate:"oneof=available on_contract sick_leave retired dead"`
	SortBy   *string `json:"sort_by" validate:"oneof=name current_rank role current_activity"`
}

type SetAdventurerHiredRequest struct {
	GuildID      uuid.UUID `json:"guild_id"`
	AdventurerID uuid.UUID `json:"adventurer_id"`
}

type ListAdventurersResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Role string    `json:"role"`
	Rank int32     `json:"rank"`
}

type ListMembersResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Role            string    `json:"role"`
	Rank            int32     `json:"rank"`
	CurrentActivity string    `json:"current_activity"`
}
type AdventurerHistoryRequest struct {
	ID       uuid.UUID `json:"adventurer_id"`
	Activity string    `json:"activity" validate:"oneof=available on_quest sick_leave retired dead"`
}

type AdventurerContractHistoryRequest struct {
	AdventurerID uuid.UUID `json:"adventurer_id"`
	ContractID   uuid.UUID `json:"contract_id"`
	Status       string    `json:"status"`
}

type AdventurersWithStatusRequest struct {
	GuildID         uuid.UUID `json:"guild_id"`
	CurrentActivity string    `json:"current_activity" validate:"oneof=available on_quest sick_leave retired dead"`
	SortBy          string    `json:"sort_by" validate:"oneof=joined_at name role"`
}

type GetMembersRequest struct {
	GuildID uuid.UUID `json:"guild_id"`
	SortBy  string    `json:"sort_by" validate:"oneof=joined_at name role activity"`
}

type GetAdventurerDetailsResponse struct {
	ID              uuid.UUID `json:"adventurer_id"`
	Name            string    `json:"name"`
	CurrentRank     int       `json:"adventurer_rank"`
	Role            string    `json:"role"`
	Description     string    `json:"description"`
	CurrentActivity string    `json:"current_activity"`
}

type AdventurerDetailsRequest struct {
	ID              uuid.UUID `json:"adventurer_id"`
	PartyID         uuid.UUID `json:"party_id"`
	Name            string    `json:"name"`
	CurrentRank     string    `json:"current_rank"`
	Role            string    `json:"role"`
	CurrentActivity string    `json:"current_activity"`
	UpkeepCost      int       `json:"upkeep_cost"`
}

type DetailsResponse struct {
	ID              uuid.UUID `json:"adventurer_id"`
	PartyID         uuid.UUID `json:"party_id"`
	Name            string    `json:"name"`
	CurrentRank     string    `json:"current_rank"`
	Role            string    `json:"role"`
	CurrentActivity string    `json:"current_activity"`
	UpkeepCost      int       `json:"upkeep_cost"`
}
