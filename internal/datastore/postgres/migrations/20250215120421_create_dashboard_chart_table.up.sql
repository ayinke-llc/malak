CREATE TYPE dashboard_chart_type AS ENUM('barchart','piechart');

CREATE TABLE dashboard_charts (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  workspace_integration_id uuid NOT NULL REFERENCES workspace_integrations(id),
  workspace_id uuid NOT NULL REFERENCES workspaces(id),
  dashboard_id uuid NOT NULL REFERENCES dashboards(id),
  dashboard_type dashboard_chart_type NOT NULL,
  reference VARCHAR (220) UNIQUE NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE dashboard_charts ADD CONSTRAINT dashboard_chart_reference_check_key CHECK (reference ~ 'dashboard_chart_[a-zA-Z0-9._]+');

ALTER TABLE dashboard_charts ADD CONSTRAINT unique_chart_per_dashboard UNIQUE(workspace_integration_id,dashboard_id,workspace_id);
