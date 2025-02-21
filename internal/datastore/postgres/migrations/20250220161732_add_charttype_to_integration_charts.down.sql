ALTER TABLE integration_charts DROP COLUMN chart_type;

ALTER TABLE dashboard_charts ADD COLUMN dashboard_type dashboard_chart_type NOT NULL; 
