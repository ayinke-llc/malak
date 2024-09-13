ALTER TABLE contacts ADD CONSTRAINT unique_email_per_workspace UNIQUE(email,workspace_id);
