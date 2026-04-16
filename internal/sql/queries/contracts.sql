-- name: ListContracts :many
SELECT
    id,
    guild_id,
    title,
    difficulty,
    minimum_party_size,
    contract_status
FROM contracts
WHERE 
    (sqlc.narg('difficulty')::int IS NULL OR difficulty = sqlc.narg('difficulty')) AND
    (sqlc.narg('min_party_size')::int IS NULL OR minimum_party_size >= sqlc.narg('min_party_size')) AND
    (sqlc.narg('status')::contract_status_enum IS NULL OR contract_status = sqlc.narg('status'))
ORDER BY
    CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'contract_status' THEN contract_status END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'minimum_party_size' THEN minimum_party_size END ASC;

-- name: GetContractByID :one
SELECT
id,
guild_id,
difficulty,
description,
minimum_party_size,
contract_status
FROM contracts
WHERE id = $1;

-- name: GetPastContractsWithStatus :many
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


-- name: SetContractStatus :exec
UPDATE contracts
SET contract_status = $2
WHERE id = $1 AND guild_id = $3;

-- name: InsertContractHistory :exec
INSERT INTO contract_history (
guild_id,
contract_id,
party_id,
status
) VALUES ($1, $2, $3, $4);

-- name: GetPartyOnContract :one
SELECT
id,
name,
party_rank
FROM parties
WHERE contract_id = $1;

-- name: AssignToGuild :exec
UPDATE contracts
SET guild_id = $2
WHERE id = $1;

-- name: CountPartyCompleteContracts :one
SELECT COUNT(*)
FROM party_history ph
JOIN contracts c ON ph.contract_id = c.id
JOIN parties p ON ph.party_id = p.id
WHERE ph.party_id = $1 
AND ph.contract_status = 'complete'
AND c.difficulty >= p.party_rank;

