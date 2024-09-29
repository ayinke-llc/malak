ALTER TABLE updates DROP COLUMN content;
ALTER TABLE updates ADD COLUMN content jsonb NOT NULL DEFAULT '{}'::jsonb;
