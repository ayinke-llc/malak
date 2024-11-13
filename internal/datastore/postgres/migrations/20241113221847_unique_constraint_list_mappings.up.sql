ALTER TABLE contact_list_mappings ADD CONSTRAINT unique_contact_list UNIQUE (contact_id, list_id);
