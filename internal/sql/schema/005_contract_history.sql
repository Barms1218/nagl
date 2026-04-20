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

CREATE INDEX ix_contract_history_guild ON contract_history (guild_id);
CREATE INDEX ix_contract_history_party ON contract_history (party_id);
CREATE INDEX ix_contract_history ON contract_history (contract_id);
CREATE INDEX ix_contract_history_diff ON contract_history (difficulty);
CREATE INDEX ix_contract_history_status ON contract_hsitory (status);

-- +goose Down

DROP INDEX ix_contract_history_guild ON contract_history (guild_id);
DROP INDEX ix_contract_history_party ON contract_history (party_id);
DROP INDEX ix_contract_history ON contract_history (contract_id);
DROP INDEX ix_contract_history_diff ON contract_history (difficulty);
DROP INDEX ix_contract_history_status ON contract_hsitory (status);
DROP TABLE contract_history;
