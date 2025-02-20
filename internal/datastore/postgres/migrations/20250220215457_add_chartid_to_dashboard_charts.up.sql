ALTER TABLE dashboard_charts DROP CONSTRAINT unique_chart_per_dashboard;

ALTER TABLE dashboard_charts ADD COLUMN chart_id uuid NOT NULL REFERENCES integration_charts(id);

ALTER TABLE dashboard_charts ADD CONSTRAINT unique_chart_per_dashboard UNIQUE(workspace_integration_id,dashboard_id,workspace_id,chart_id);
