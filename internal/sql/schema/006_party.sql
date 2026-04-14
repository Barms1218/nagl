-- +goose Up
CREATE TABLE parties (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	guild_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
	contract_id UUID NOT NULL REFERENCES contracts(id) ON DELETE CASCADE
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	UNIQUE(guild_id, contract_id)
);

-- +goose Down
DROP TABLE parties;
