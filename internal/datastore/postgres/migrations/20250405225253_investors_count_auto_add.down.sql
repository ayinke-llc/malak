DROP TRIGGER IF EXISTS increment_investors_count_on_contact_insert ON fundraising_pipeline_column_contacts;
DROP TRIGGER IF EXISTS decrement_investors_count_on_contact_delete ON fundraising_pipeline_column_contacts;
DROP FUNCTION IF EXISTS increment_investors_count();
DROP FUNCTION IF EXISTS decrement_investors_count();
