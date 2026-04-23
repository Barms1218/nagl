package procedural

import (
	"github.com/google/uuid"
)

type GeneratedContract struct {
	Title        string `json:"title" jsonschema:"description=Fitting name for contract"`
	Difficulty   int32  `json:"difficulty" jsonschema:"description=1-5 denoting task difficulty"`
	RecPartySize int    `json:"rec_party_size" jsonschema:"description=1-5 taking into account difficulty"`
	Description  string `json:"description" jsonschema:"description=3-4 sentences of flavor text"`
	Reward       int    `json:"reward" jsonschema:"description=Reward for completion"`
	Duration     int    `json:"duration" jsonschema:"description=Duration in minutes"`
}

type GeneratedAdventurer struct {
	Name            string `json:"name" jsonschema:"description=Full name"`
	Role            string `json:"role" jsonschema:"description=Frontliner | Spellcaster | Healer | Generalist"`
	Rank            int32  `json:"rank" jsonschema:"description=1-5 integer denoting rank"`
	Bio             string `json:"bio" jsonschema:"description=2-3 sentences of flavor text"`
	UpkeepCost      int    `json:"upkeep_cost" jsonschema:"description=Cost per day"`
	RecruitmentCost int    `json:"recruitment_cost" jsonschema:"descriptino=Cost to recruit"`
}

type ContractRequest struct {
	Title        string `json:"title"`
	Difficulty   int    `json:"difficulty"`
	RecPartySize int    `json:"rec_party_size"`
	Description  string `json:"description"`
	Reward       int    `json:"reward"`
	Duration     int    `json:"duration" validate:"gte=60,lte=360"`
}

type GeneratedParty struct {
	GuildID      uuid.UUID `json:"guild_id" jsonschema:"description=The guid ID"`
	GuildName    string    `json:"guild_name" jsonschema:"description=The guild name"`
	PartyName    string    `json:"party_name" jsonschema:"description=The party name"`
	MaxPartySize int32     `json:"max_party_size"`
	PartyStatus  string    `json:"party_status"`
}

type PartyName struct {
	GuildID   uuid.UUID `json:"guild_id" jsonschema:"description=The guid ID"`
	PartyName string    `json:"party_name" jsonschema:"description=The party name"`
}

type GeneratePartyRequest struct {
	GuildID     uuid.UUID     `json:"guild_id"`
	GuildName   string        `json:"guild_name"`
	Adventurers []PartyMember `json:"adventurers"`
}

type PartyMember struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Rank int    `json:"rank"`
}
