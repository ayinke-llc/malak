
CREATE TABLE integration_charts (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  workspace_id uuid NOT NULL REFERENCES workspaces(id),
  workspace_integration_id uuid NOT NULL REFERENCES workspace_integrations(id),
  reference VARCHAR (220) UNIQUE NOT NULL,
  user_facing_name VARCHAR(220) NOT NULL,
  internal_name VARCHAR(220) NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE integration_charts ADD CONSTRAINT integration_charts_reference_check_key CHECK (reference ~ 'integration_chart_[a-zA-Z0-9._]+');

ALTER TABLE integration_charts ADD CONSTRAINT unique_internal_name_per_workspace_and_integration UNIQUE(internal_name,workspace_id,workspace_integration_id);

ALTER TABLE integration_data_points ADD integration_chart_id uuid REFERENCES integration_charts(id);
