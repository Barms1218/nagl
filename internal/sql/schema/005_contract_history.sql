-- +goose Up
CREATE TABLE contract_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uui(),
	guild_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
	contract_id UUID NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status contract_status NOT NULL	
	UNIQUE(guild_id, contract_id)
);
