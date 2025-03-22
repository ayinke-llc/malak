ALTER TABLE integration_charts DROP COLUMN data_point_type;

ALTER TABLE integration_data_points ADD COLUMN data_point_type integration_data_point_type NOT NULL;
