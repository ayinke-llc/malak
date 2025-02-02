CREATE TYPE integration_data_point_type AS ENUM('currency', 'others');

CREATE TABLE integration_data_points (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  workspace_id uuid NOT NULL REFERENCES workspaces(id),
  workspace_integration_id uuid NOT NULL REFERENCES workspace_integrations(id),
  data_point_type integration_data_point_type NOT NULL,
  reference VARCHAR (220) UNIQUE NOT NULL,
  point_name VARCHAR(200) NOT NULL,
  point_value BIGINT NOT NULL,

  metadata jsonb NOT NULL DEFAULT '{}'::jsonb,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE integration_data_points ADD CONSTRAINT integration_data_point_reference_check_key CHECK (reference ~ 'integration_datapoint_[a-zA-Z0-9._]+');
