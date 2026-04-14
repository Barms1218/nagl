-- name: GetContracts :many
SELECT
id,
guild_id,
title,
difficulty,
minimum_party_size,
contract_status
FROM contracts
ORDER BY
CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
CASE WHEN sqlc.arg(sort_by)::text = 'contract_status' THEN contract_status END ASC,
CASE WHEN sqlc.arg(sort_by)::text = 
	'minimum_party_size' THEN minimum_party_size END ASC;

-- name: GetContractsWithDifficulty :many
SELECT
id,
guild_id,
title,
difficulty,
minimum_party_size,
contract_status
FROM contracts
WHERE difficulty = $1
ORDER BY
CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
CASE WHEN sqlc.arg(sort_by)::text = 'contract_status' THEN contract_status END ASC,
CASE WHEN sqlc.arg(sort_by)::text = 
	'minimum_party_size' THEN minimum_party_size END ASC;

-- name: GetContractsWithStatus :many
SELECT 
ch.id AS contract_history_id,
ch.guild_id,
ch.contract_id,
c.title,
ch.occurred_at,
ch.status,
c.difficulty
FROM contract_history ch
JOIN contracts c ON ch.contract_id = c.id
WHERE ch.status = $1
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 
		'minimum_party_size' THEN minimum_party_size END ASC;

-- name: GetContractsWithMinPartySize :many
SELECT
id,
guild_id,
title,
difficulty,
minimum_party_size,
contract_status
FROM contracts
WHERE minimum_party_size = $1
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'contract_status' THEN contract_status END ASC;

-- name: InsertContract :one
INSERT INTO contracts (
	guild_id,
	title, 
	difficulty,
	description,
	minimum_party_size
)VALUES($1, $2, $3, $4, $5)
RETURNING *;

-- name: SetContractStatus :exec
UPDATE contracts
SET party_id = NULL, contract_status = $2
WHERE id = $1;

-- name: InsertContractHistory :exec
INSERT INTO contract_history (
guild_id,
contract_id,
status
) VALUES ($1, $2, $3);

-- name: GetPartyOnContract :many
SELECT
p.id AS party_id,
a.id AS adventurer_id,
a.name,
a.current_rank,
a.role
FROM contracts c
JOIN parties p ON c.id = p.contract_id
JOIN party_members pm ON p.id = pm.party_id
JOIN adventurers a ON pm.adventurer_id = a.id
WHERE c.id = $1;
