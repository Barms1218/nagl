-- +goose Up
CREATE TABLE parties (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	contract_id UUID NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
	name TEXT UNIQUE NOT NULL,
	guild_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
	maximum_party_size INTEGER NOT NULL DEFAULT 5,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(guild_id, contract_id)
);

-- +goose Down
DROP TABLE parties;
