ALTER TABLE contacts
    ALTER COLUMN created_by DROP NOT NULL,
    ALTER COLUMN owner_id DROP NOT NULL;
