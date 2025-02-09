CREATE TYPE integration_sync_checkpoint_status AS ENUM ('failed', 'success', 'pending');

CREATE TABLE integration_sync_checkpoints (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  workspace_id uuid NOT NULL REFERENCES workspaces(id),
  workspace_integration_id uuid NOT NULL REFERENCES workspace_integrations(id),
  last_sync_attempt TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  last_successful_sync TIMESTAMP WITH TIME ZONE,
  status integration_sync_checkpoint_status NOT NULL DEFAULT 'pending',
  error_message TEXT NULL DEFAULT '',
  reference VARCHAR (220) UNIQUE NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE integration_sync_checkpoints 
  ADD CONSTRAINT integration_sync_checkpoint_reference_check_key CHECK (reference ~ 'integration_sync_checkpoint_[a-zA-Z0-9._]+');
