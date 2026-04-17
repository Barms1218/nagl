package guild

type GuildAuthRequest struct {
	GuildName string `json:"guild_name"`
	GuildKey  string `json:"guild_key" validate:"min=8,max=13"`
}
