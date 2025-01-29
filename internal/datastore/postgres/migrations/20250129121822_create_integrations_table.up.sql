CREATE TYPE integration_connection_type AS ENUM('oauth2', 'api_key', 'malak');

CREATE TABLE integrations (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  integration_name TEXT UNIQUE NOT NULL,
  reference VARCHAR (220) UNIQUE NOT NULL,
  description TEXT NOT NULL,
  is_enabled BOOLEAN DEFAULT TRUE ,
  integration_type integration_connection_type NOT NULL,

  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE integrations ADD CONSTRAINT integration_reference_check_key CHECK (reference ~ 'integration_[a-zA-Z0-9._]+');

CREATE TABLE workspace_integrations (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,
  workspace_id uuid NOT NULL REFERENCES workspaces(id),
  integration_id uuid NOT NULL REFERENCES integrations(id),
  is_enabled BOOLEAN DEFAULT TRUE ,

  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE workspace_integrations ADD CONSTRAINT workspace_integration_reference_check_key CHECK (reference ~ 'workspace_integration_[a-zA-Z0-9._]+');
