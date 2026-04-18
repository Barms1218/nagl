-- name: ListAvailableContracts :many
SELECT
id,
title,
difficulty,
minimum_party_size,
contract_status
FROM contracts
WHERE guild_id IS NULL
ORDER BY
    CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'contract_status' THEN contract_status END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'minimum_party_size' THEN minimum_party_size END ASC;

-- name: ListContractsWithMinDifficulty :many
SELECT
id,
title,
difficulty,
minimum_party_size,
contract_status
FROM contracts
WHERE difficulty >= $1 AND guild_id = $2
ORDER BY
    CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'contract_status' THEN contract_status END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'minimum_party_size' THEN minimum_party_size END ASC;

-- name: ListContractsWithMaxDifficulty :many
SELECT
id,
title,
difficulty,
minimum_party_size,
contract_status
FROM contracts
WHERE difficulty <= $1 AND guild_id = $2
ORDER BY
    CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'contract_status' THEN contract_status END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'minimum_party_size' THEN minimum_party_size END ASC;

-- name: ListContractsWithStatus :many
SELECT
id,
title,
difficulty,
minimum_party_size,
contract_status
FROM contracts
WHERE contract_status = $1 AND guild_id = $2
ORDER BY
    CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'minimum_party_size' THEN minimum_party_size END ASC;

-- name: ListContractsWithMinPartySize :many
SELECT
id,
title,
difficulty,
minimum_party_size,
contract_status
FROM contracts
WHERE minimum_party_size >= $1 AND guild_id = $2
ORDER BY
    CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'contract_status' THEN contract_status END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'minimum_party_size' THEN minimum_party_size END ASC;

-- name: ListContractsWithMaxPartySize :many
SELECT
id,
title,
difficulty,
minimum_party_size,
contract_status
FROM contracts
WHERE minimum_party_size <= $1 AND guild_id = $2
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

