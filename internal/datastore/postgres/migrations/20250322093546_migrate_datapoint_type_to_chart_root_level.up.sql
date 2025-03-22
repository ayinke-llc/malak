ALTER TABLE integration_data_points DROP COLUMN data_point_type;

ALTER TABLE integration_charts ADD COLUMN data_point_type integration_data_point_type NOT NULL DEFAULT 'others';
