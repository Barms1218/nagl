-- +goose Up
CREATE FUNCTION party_size()
RETURNS TRIGGER AS $$
DECLARE
	new_size INTEGER;
BEGIN
	SELECT COUNT(*) into new_size
	FROM party_members
	WHERE party_id = NEW.party_id;

	IF current_size > 5
		RAISE EXCEPTION 'Party cannot exceed 5 members'
	END if;

	RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER check_party_size
BEFORE INSERT ON party_members
FOR EACH ROW
EXECUTE FUNCTION party_size();

-- +goose Down
DROP FUNCTION party_size();
DROP TRIGGER check_party_size();
