-- name: CreateParty :one
INSERT INTO parties (guild_id, contract_id)
VALUES($1, $2)RETURNING id;

-- name: GetAllParties :many
SELECT 

-- name: GetPartyDetails :many
SELECT 
p.id AS party_id,
c.title,
c.contract_status,
a.id AS adventurer_id,
a.current_activity,
a.name,
a.role
FROM parties p
JOIN contracts c ON p.contract_id = c.id
JOIN party_members pm ON pm.party_id = p.id
JOIN adventurers a ON a.id = pm.adventurer_id
WHERE p.id = $1

