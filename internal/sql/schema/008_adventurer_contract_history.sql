-- +goose Up
CREATE TABLE adventurer_contract_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	adventurer_id UUID NOT NULL REFERENCES adventurers(id) ON DELETE CASCADE,
	contract_id UUID NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
	status contract_status_enum NOT NULL,
	occurred_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX ix_adventurer_id ON adventurer_contract_history (adventurer_id);
CREATE INDEX ix_contract_id ON adventurer_contract_history (contract_id);
CREATE INDEX ix_conract_status ON adventurer_contract_history (status);

-- +goose Down
DROP INDEX ix_adventurer_id ON adventurer_contract_history (adventurer_id);
DROP INDEX ix_contract_id ON adventurer_contract_history (contract_id);
DROP INDEX ix_conract_status ON adventurer_contract_history (status);
DROP TABLE adventurer_contract_history;
