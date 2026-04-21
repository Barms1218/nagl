package procedural

import (
	"github.com/google/uuid"
)

type GeneratedContract struct {
	GuildID      uuid.UUID `json:"guild_id"`
	Title        string    `json:"title"`
	Difficulty   string    `json:"difficulty"`
	RecPartySize int       `json:"rec_party_size"`
	Description  string    `json:"description"`
	Reward       int       `json:"reward"`
}

type AdventurerRequest struct {
	Name            string `json:"name"`
	Role            string `json:"role"`
	Rank            int    `json:"rank"`
	Bio             string `json:"bio"`
	UpkeepCost      int    `json:"upkeep_cost" validate:"gte=25, lte=50"`
	RecruitmentCost int    `json:"recruitment_cost" validate:"gte=50,lte=500"`
}

type GeneratedAdventurer struct {
	Name            string `json:"name"`
	Role            string `json:"role"`
	Rank            string `json:"rank"`
	Bio             string `json:"bio"`
	UpkeepCost      int    `json:"upkeep_cost" validate:"gte=25, lte=50"`
	RecruitmentCost int    `json:"recruitment_cost" validate:"gte=50,lte=500"`
}

type ContractRequest struct {
	Title        string `json:"title"`
	Difficulty   int    `json:"difficulty"`
	RecPartySize int    `json:"rec_party_size"`
	Description  string `json:"description"`
	Reward       int    `json:"reward"`
}
