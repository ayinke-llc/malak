ALTER TABLE dashboards ADD COLUMN workspace_id uuid NOT NULL REFERENCES workspaces(id);
