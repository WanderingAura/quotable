CREATE FUNCTION sync_last_modified() RETURNS trigger as $$
BEGIN
    NEW.last_modified := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER quotes_modified_trigger BEFORE UPDATE ON quotes
    FOR EACH ROW EXECUTE PROCEDURE sync_last_modified(); 