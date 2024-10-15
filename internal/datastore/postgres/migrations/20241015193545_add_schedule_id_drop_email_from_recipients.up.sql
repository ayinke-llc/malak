ALTER TABLE update_recipients DROP COLUMN email;
ALTER TABLE update_recipients DROP COLUMN recipient_type;

ALTER TABLE update_recipients ADD COLUMN contact_id uuid NOT NULL REFERENCES contacts(id);
ALTER TABLE update_recipients ADD COLUMN schedule_id uuid NOT NULL REFERENCES update_schedules(id);
