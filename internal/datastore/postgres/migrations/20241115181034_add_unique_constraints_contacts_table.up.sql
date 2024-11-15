ALTER TABLE contacts ADD CONSTRAINT unique_contact_workspace UNIQUE (email, workspace_id);
