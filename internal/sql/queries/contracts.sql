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
CASE WHEN sqlc.arg(sort_by)::text = 'contract_status' THEN contract_status END ASC
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
id,
guild_id,
title,
difficulty,
minimum_party_size,
contract_status
FROM contracts
WHERE contract_status = $1;
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
CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
CASE WHEN sqlc.arg(sort_by)::text = 'contract_status' THEN contract_status END ASC;

-- name: UpsertContracts :batchone
INSERT INTO contracts (guild_id, title, difficulty, minimum_party_size)
VALUES($1, $2, $3, $4)
RETURNING id, title, difficulty, minimum_party_size, contract_statusl
