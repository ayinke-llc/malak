DROP TRIGGER IF EXISTS update_fundraising_pipeline_columns_updated_at ON fundraising_pipeline_columns;
DROP TRIGGER IF EXISTS update_fundraising_pipelines_updated_at ON fundraising_pipelines;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS fundraising_pipeline_columns;
DROP TABLE IF EXISTS fundraising_pipelines;

DROP TYPE IF EXISTS fundraise_pipeline_column_type;
DROP TYPE IF EXISTS fundraise_pipeline_stage;
