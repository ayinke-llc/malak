CREATE TYPE update_recipients_status AS ENUM ('pending', 'sent', 'failed');

ALTER TABLE update_recipients ADD COLUMN status update_recipients_status DEFAULT 'pending';
