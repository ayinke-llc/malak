CREATE OR REPLACE FUNCTION increment_investors_count()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE fundraising_pipeline_columns
    SET investors_count = investors_count + 1
    WHERE id = NEW.fundraising_pipeline_column_id;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER increment_investors_count_on_contact_insert
    AFTER INSERT ON fundraising_pipeline_column_contacts
    FOR EACH ROW
    WHEN (NEW.deleted_at IS NULL)
    EXECUTE FUNCTION increment_investors_count();

-- Also add a trigger for when contacts are soft deleted to decrement the count
CREATE OR REPLACE FUNCTION decrement_investors_count()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.deleted_at IS NOT NULL AND OLD.deleted_at IS NULL THEN
        UPDATE fundraising_pipeline_columns
        SET investors_count = investors_count - 1
        WHERE id = NEW.fundraising_pipeline_column_id;
    END IF;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER decrement_investors_count_on_contact_delete
    BEFORE UPDATE ON fundraising_pipeline_column_contacts
    FOR EACH ROW
    EXECUTE FUNCTION decrement_investors_count();
