CREATE TABLE contact_lists (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  workspace_id uuid NOT NULL REFERENCES workspaces(id),
  title VARCHAR (200) NOT NULL,
  reference VARCHAR (220) UNIQUE NOT NULL,
  created_by uuid NOT NULL REFERENCES users(id),

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE contact_lists ADD CONSTRAINT contact_list_reference_check_key CHECK (reference ~ 'list_[a-zA-Z0-9._]+');
