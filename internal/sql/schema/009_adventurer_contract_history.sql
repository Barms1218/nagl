-- +goose Up
CREATE TABLE adventurer_contract_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	adventurer_id UUID NOT NULL REFERENCES adventurers(id) ON DELETE CASCADE,
	contract_id UUID NOT NULL REFERENCES contracts(id) ON DELETE CASCADE,
	status contract_status NOT NULL,
	occurred_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE adventurer_contract_history;
