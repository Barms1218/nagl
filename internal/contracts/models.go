package contracts

import (
	"github.com/google/uuid"
)

// Search parameters for a contract
type SearchFilters struct {
	MinDifficulty *int32  `json:"min_difficulty" validate:"omitempty,gte=1,lte=5"`
	MaxDifficulty *int32  `json:"max_difficulty" validate:"omitempty,gte=1,lte=5"`
	PartySize     *int32  `json:"partySize" validate:"omitempty,gte=1,lte=5"`
	Status        *string `json:"status" validate:"omitempty,oneof=complete failed in-progress available"`
	SortBy        *string `json:"sortBy" validate:"oneof=title difficulty minimum_party_size contract_status"`
}

type ListContractsResponse struct {
	ID             uuid.UUID `json:"contract_id"`
	Title          string    `json:"title"`
	Difficulty     int32     `json:"difficulty"`
	RecPartySize   int32     `json:"minimum_party_size" validate:"gte=1,lte=5"`
	ContractStatus string    `json:"contract_status"`
}

type ActiveContractDetailsResponse struct {
	ID           uuid.UUID `json:"contract_id"`
	Title        string    `json:"title"`
	GuildName    string    `json:"guild_name"`
	PartyName    string    `json:"party_name"`
	PartyStatus  string    `json:"party_status"`
	Description  string    `json:"description"`
	Difficulty   int32     `json:"difficulty"`
	RecPartySize int32     `json:"rec_party_size"`
	Status       string    `json:"status"`
}

type ContractDetailsResponse struct {
	ID           uuid.UUID `json:"contract_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Difficulty   int32     `json:"difficulty"`
	RecPartySize int32     `json:"rec_party_size"`
	Status       string    `json:"status"`
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

type ContractWithStatusRequest struct {
	ContractStatus string    `json:"contract_status" validate:"oneof=complete failed in-progress available"`
	GuildID        uuid.UUID `json:"guild_id"`
	SortBy         string    `json:"sort_by" validate:"oneof=title difficulty minimum_party_size contract_status"`
}

type ContractWithDifficultyRequeset struct {
	Difficulty int32     `json:"difficulty" validate:"gte=1,lte=5"`
	GuildID    uuid.UUID `json:"guild_id"`
	SortBy     string    `json:"sort_by" validate:"oneof=title difficulty minimum_party_size contract_status"`
}
