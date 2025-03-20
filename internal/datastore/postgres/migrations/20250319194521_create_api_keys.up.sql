CREATE TABLE api_keys (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  workspace_id uuid NOT NULL REFERENCES workspaces(id),
  created_by uuid NOT NULL REFERENCES users(id),
  reference VARCHAR (220) UNIQUE NOT NULL,
  value VARCHAR (220) UNIQUE NOT NULL,
  key_name VARCHAR(220) NOT NULL,

  expires_at TIMESTAMP WITH TIME ZONE,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE api_keys ADD CONSTRAINT reference_check_key 
  CHECK (reference ~ 'api_key_[a-zA-Z0-9._]+');
