CREATE TYPE chart_type AS ENUM('bar', 'pie');

ALTER TABLE integration_charts ADD COLUMN chart_type chart_type NOT NULL;

ALTER TABLE dashboard_charts DROP COLUMN dashboard_type; 
