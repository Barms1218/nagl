-- +goose Up
CREATE TABLE party_members (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	adventurer_id UUID NOT NULL REFERENCES adventurers(id) ON DELETE CASCADE,
	party_id UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
	maximum_party_size INTEGER NOT NULL 5
	UNIQUE(adventurer_id, party_id)
);

-- +goose Down
DROP TABLE party_members;
