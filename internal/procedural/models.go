package procedural

import (
	"github.com/google/uuid"
)

type GeneratedContract struct {
	GuildID          uuid.UUID `json:"guild_id"`
	Title            string    `json:"title"`
	Difficulty       int       `json:"difficulty"`
	MinimumPartySize int       `json:"min_party_size"`
	Description      string    `json:"description"`
	Reward           int       `json:"reward"`
}
