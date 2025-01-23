ALTER TABLE update_recipients ADD CONSTRAINT unique_contact_per_update UNIQUE(contact_id,update_id);
