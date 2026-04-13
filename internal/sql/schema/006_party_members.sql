-- +goose Up
CREATE TABLE party_members (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	party_id UUID NOT NULL REFERENCES parties(id) ON DELETE CASCADE,
)
