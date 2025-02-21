CREATE TABLE dashboard_chart_positions (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  dashboard_id uuid NOT NULL REFERENCES dashboards(id),
  chart_id uuid NOT NULL REFERENCES dashboard_charts(id),
  order_index INT NOT NULL
);

ALTER TABLE dashboard_chart_positions ADD CONSTRAINT one_position_per_chart UNIQUE (dashboard_id, chart_id, order_index); 
