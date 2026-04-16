-- +goose Up
CREATE TABLE contract_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	guild_id UUID NOT NULL REFERENCES guilds(id) ON DELETE CASCADE,
	party_id UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
	contract_id UUID NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
	difficulty INTEGER NOT NULL DEFAULT 1,
	occurred_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	status contract_status_enum NOT NULL
);

-- +goose Down
DROP TABLE contract_history;
