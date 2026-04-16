package contracts

import (
	"github.com/google/uuid"
)

// Search parameters for a contract
type SearchFilters struct {
	Difficulty *int    `json:"difficulty" validate:"omitempty,gte=1,lte=5"`
	PartySize  *int    `json:"partySize" validate:"omitempty,gte=1,lte=5"`
	Status     *string `json:"status" validate:"omitempty,oneof=complete failed in-progress available"`
	SortBy     string  `json:"sortBy" validate:"oneof=title difficulty minimum_party_size contract_status"`
}

// Parameters for the history of a contract
type HistoryParams struct {
	GuildID    uuid.UUID `json:"guild_id"`
	ContractID uuid.UUID `json:"contract_id"`
	Status     string    `json:"status"`
}

// Parameters for finding past contracts with
// Failed/Complete Status
type PastContractsParams struct {
	Status string `json:"status" validate:"oneof= complete failed"`
	SortBy string `json:"sortBy" validate:"oneof=title difficulty minimum_party_size contract_status"`
}

// Parameters to set a particular contract to a particular status
type SetContractStatusRequest struct {
	GuildID   uuid.UUID `json:"guild_id"`
	ID        uuid.UUID `json:"id"`
	PartyID   uuid.UUID `json:"party_id"`
	NewStatus string    `json:"status" validate:"omitempty,oneof=complete failed in-progress available"`
}

type ContractClaimRequest struct {
	ContractID uuid.UUID `json:"contract_id"`
	GuildID    uuid.UUID `json:"guild_id"`
}
