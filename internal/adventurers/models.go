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

type GetMembersRequest struct {
	GuildID uuid.UUID `json:"guild_id"`
	SortBy  string    `json:"sort_by" validate:"oneof=joined_at name role activity"`
}

type GetMembersResponse struct {
	ID              uuid.UUID `json:"adventurer_id"`
	CurrentRank     int       `json:"adventurer_rank"`
	CurrentActivity string    `json:"current_activity"`
	Name            string    `json:"name"`
	Role            string    `json:"role"`
}

type AdventurerDetailsRequest struct {
	PartyID         uuid.UUID `json:"party_id"`
	Name            string    `json:"name"`
	CurrentRank     int       `json:"current_rank"`
	Role            string    `json:"role"`
	CurrentActivity string    `json:"current_activity"`
	UpkeepCost      int       `json:"upkeep_cost"`
}

type DetailsResponse struct {
	PartyID         uuid.UUID `json:"party_id"`
	Name            string    `json:"name"`
	CurrentRank     int       `json:"current_rank"`
	Role            string    `json:"role"`
	CurrentActivity string    `json:"current_activity"`
	UpkeepCost      int       `json:"upkeep_cost"`
}
