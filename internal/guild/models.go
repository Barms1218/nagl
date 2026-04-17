package guild

import "github.com/google/uuid"

type GuildAuthRequest struct {
	GuildName string `json:"guild_name"`
	GuildKey  string `json:"guild_key" validate:"min=8,max=13"`
}

type UpdateTreasuryRequest struct {
	GuildID uuid.UUID `json:"guild_id"`
	Amount  int32     `json:"amount"`
}
