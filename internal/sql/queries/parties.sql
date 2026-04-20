-- name: CreateParty :one
INSERT INTO parties (guild_id, contract_id, name)
VALUES($1, $2, $3)
RETURNING *;

-- name: GetParty :one
SELECT
id,
guild_id,
contract_id,
party_rank,
name,
maximum_party_size
FROM parties
WHERE id = $1;

-- name: GetAllParties :many
SELECT 
p.id AS party_id,
c.id AS contract_id,
c.title,
c.contract_status,
a.name,
a.role,
a.current_activity
FROM parties p
JOIN contracts c ON p.contract_id = c.id
JOIN adventurers a ON p.id = a.party_id
WHERE p.guild_id = $1;

-- name: GetPartyDetails :many
SELECT 
p.id AS party_id,
c.title,
c.contract_status,
p.name,
p.party_rank,
a.id AS adventurer_id,
a.current_activity,
a.name,
a.role
FROM parties p
JOIN contracts c ON p.contract_id = c.id
JOIN adventurers a ON p.id = a.party_id
WHERE p.id = $1;

-- name: InsertPartyHistory :exec
INSERT INTO party_history (
	party_id,
    contract_id,
	contract_status
) VALUES($1 , $2, $3);

-- name: SetMemberStatus :exec
UPDATE adventurers
SET current_activity = $1
WHERE party_id = (
    SELECT id 
    FROM parties 
    WHERE contract_id = $2
);

-- name: SetPartyRank :exec
UPDATE parties
SET party_rank = $1
WHERE id = $2;

-- name: CountMemberCompleteContracts :many
SELECT ach.adventurer_id, a.name, a.current_rank, COUNT(*) AS completed_count
FROM adventurer_contract_history ach
JOIN adventurers a ON a.id = ach.adventurer_id
JOIN contracts c ON ach.contract_id = c.id
WHERE ach.contract_id = $1
AND c.difficulty >= a.current_rank
AND ach.status = 'complete'
GROUP BY ach.adventurer_id, a.name, a.current_rank;


-- name: InsertMemberContractHistory :exec
INSERT INTO adventurer_contract_history(
adventurer_id,
contract_id,
status
)
SELECT
a.id,
c.id,
$2
FROM adventurers a
JOIN parties p ON a.party_id = p.id
WHERE p.contract_id = $1;

-- name: RemovePartyFromContract :exec
UPDATE parties
SET contract_id = NULL
WHERE id = $1;

-- name: AddpartyToContract :exec
Update parties
SET contract_id = $2
WHERE id = $1;
