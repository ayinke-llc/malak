ALTER TABLE integration_charts ADD COLUMN metadata jsonb NOT NULL DEFAULT '{}'::jsonb;
