package models

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
