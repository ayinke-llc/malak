ALTER TABLE update_recipient_stats ADD COLUMN update_id UUID NOT NULL REFERENCES updates(id);
