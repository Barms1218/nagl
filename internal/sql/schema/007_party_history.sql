-- +goose Up
CREATE TABLE party_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	party_id UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
	contract_id UUID NOT NULL REFERENCES contracts(id) ON DELETE SET NULL,
	occurred_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	contract_status contract_status_enum NOT NULL
);

CREATE INDEX ix_party_id ON party_history (party_id);
CREATE INDEX ix_party_contract ON party_history (contract_id);
CREATE INDEX ix_contract_status ON party_history (contract_status);

-- +goose Down
DROP INDEX ix_party_id ON party_history (party_id);
DROP INDEX ix_party_contract ON party_history (contract_id);
DROP INDEX ix_contract_status ON party_history (contract_status);
DROP TABLE party_history;
