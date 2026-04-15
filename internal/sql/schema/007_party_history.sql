-- +goose Up
CREATE TABLE party_history (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	party_id UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
	occured_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	contract_status contract_status_enum NOT NULL
);

-- +goose Down
DROP TABLE party_members;
