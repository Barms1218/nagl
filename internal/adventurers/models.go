package adventurers

import (
	"github.com/google/uuid"
)

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

type GetAdventurersResponse struct {
	ID              uuid.UUID `json:"adventurer_id"`
	Name            string    `json:"name"`
	CurrentRank     int       `json:"adventurer_rank"`
	Role            string    `json:"role"`
	CurrentActivity string    `json:"current_activity"`
}

type AdventurerDetailsRequest struct {
	ID              uuid.UUID `json:"adventurer_id"`
	PartyID         uuid.UUID `json:"party_id"`
	Name            string    `json:"name"`
	CurrentRank     int       `json:"current_rank"`
	Role            string    `json:"role"`
	CurrentActivity string    `json:"current_activity"`
	UpkeepCost      int       `json:"upkeep_cost"`
}

type DetailsResponse struct {
	ID              uuid.UUID `json:"adventurer_id"`
	PartyID         uuid.UUID `json:"party_id"`
	Name            string    `json:"name"`
	CurrentRank     int       `json:"current_rank"`
	Role            string    `json:"role"`
	CurrentActivity string    `json:"current_activity"`
	UpkeepCost      int       `json:"upkeep_cost"`
}
