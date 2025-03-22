ALTER TABLE integration_data_points
ADD CONSTRAINT unique_integration_data_point 
UNIQUE(workspace_id, workspace_integration_id, integration_chart_id, point_name); 
