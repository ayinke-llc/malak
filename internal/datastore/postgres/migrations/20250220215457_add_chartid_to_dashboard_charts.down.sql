ALTER TABLE dashboard_charts ADD CONSTRAINT unique_chart_per_dashboard UNIQUE(workspace_integration_id,dashboard_id,workspace_id);
