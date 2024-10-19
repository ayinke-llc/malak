CREATE TYPE recipient_type_enum AS ENUM ('list', 'email');

ALTER TABLE update_recipients ADD COLUMN recipient_type recipient_type_enum NOT NULL;
