-- name: InsertContractHistory :exec
INSERT INTO contract_history (
guild_id,
contract_id,
party_id,
difficulty,
status
) VALUES ($1, $2, $3, $4, $5);

-- name: GetHistoryOfContract :many
SELECT
ch.id AS contract_history_id,
g.name AS guild_name,
p.name AS party_name,
ch.contract_id,
c.title,
ch.occurred_at,
ch.status,
ch.difficulty
FROM contract_history ch
JOIN contracts c ON ch.contract_id = c.id
JOIN guilds g ON ch.guild_id = g.id
JOIN parties p ON ch.party_id = p.id
WHERE ch.contract_id = $1 AND ch.guild_id = $2;

-- name: GetPartyContractHistory :many
SELECT 
ch.id AS contract_history_id,
c.title,
p.name,
p.party_rank,
ch.occurred_at,
ch.difficulty,
ch.status
FROM contract_history ch
JOIN contracts c ON ch.contract_id = c.id
JOIN parties p ON ch.party_id = p.id
WHERE p.id = $1 AND ch.guild_id = $2;

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
WHERE ch.status = $1 AND ch.guild_id = $2
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 
		'rec_party_size' THEN minimum_party_size END ASC;

-- name: GetPastContractsMinDifficulty :many
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
WHERE ch.difficulty >= $1 AND ch.guild_id = $2
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 
		'rec:w
		_party_size' THEN minimum_party_size END ASC;

-- name: GetPastContractsMaxDifficulty :many
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
WHERE ch.difficulty <= $1 AND ch.guild_id = $2
ORDER BY
	CASE WHEN sqlc.arg(sort_by)::text = 'title' THEN title END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 'difficulty' THEN difficulty END ASC,
	CASE WHEN sqlc.arg(sort_by)::text = 
		'rec_party_size' THEN minimum_party_size END ASC;


