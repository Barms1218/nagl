-- +goose Up
CREATE TYPE role_enum AS ENUM('leader', 'specialist', 'scount')
CREATE TABLE parties (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	contract_id UUID NOT NULL REFERENCES contracts(id) ON DELETE CASCADE
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
