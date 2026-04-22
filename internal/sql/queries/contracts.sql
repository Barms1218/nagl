-- name: ListAvailableContracts :many
SELECT
id,
title,
difficulty,
rec_party_size,
contract_status
FROM contracts
WHERE guild_id IS NULL
    AND (sqlc.narg('max_difficulty')::int IS NULL OR difficulty <= sqlc.narg('max_difficulty'))
    AND (sqlc.narg('min_difficulty')::int IS NULL OR difficulty >= sqlc.narg('min_difficulty'))
    AND (sqlc.narg('status')::contract_status_enum IS NULL OR contract_status = sqlc.narg('status'))
    AND (sqlc.narg('party_size')::int IS NULL OR rec_party_size = sqlc.narg('party_size'))
ORDER BY
    CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'contract_status' THEN contract_status END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'rec_party_size' THEN rec_party_size END ASC;

-- name: ListGuildContracts :many
SELECT
    id,
    title,
    difficulty,
    rec_party_size,
    contract_status
FROM contracts
WHERE guild_id = $1
    AND (sqlc.narg('max_difficulty')::int IS NULL OR difficulty <= sqlc.narg('max_difficulty'))
    AND(sqlc.narg('min_difficulty')::int IS NULL OR difficulty >= sqlc.narg('min_difficulty'))
    AND (sqlc.narg('status')::contract_status_enum IS NULL OR contract_status = sqlc.narg('status'))
    AND (sqlc.narg('party_size')::int IS NULL OR rec_party_size = sqlc.narg('party_size'))
ORDER BY
    CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'contract_status' THEN contract_status END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'rec_party_size' THEN rec_party_size END ASC;

-- name: GetContractDetailsByID :one
SELECT
c.id,
c.title,
g.name AS guild_name,
p.name AS party_name,
p.party_status,
c.difficulty,
c.description,
c.rec_party_size,
c.contract_status
FROM contracts c
JOIN guilds g ON c.guild_id = g.id
JOIN parties p ON p.contract_id = c.id
WHERE c.id = $1;


-- name: GetAvailableContractDetails :one
SELECT 
id,
title,
difficulty,
description,
rec_party_size,
contract_status
FROM contracts
WHERE id = $1;


-- name: SetContractStatus :exec
UPDATE contracts
SET contract_status = $2
WHERE id = $1 AND guild_id = $3;



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

-- name: InsertContract :one
INSERT INTO contracts (
    title,
    difficulty,
    rec_party_size,
    description,
    reward
) VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAdventurersOnContract :many
SELECT
a.id
FROM adventurers a
JOIN parties p ON a.party_id = p.id
WHERE p.contract_id = $1;
