CREATE TYPE fundraise_pipeline_stage AS ENUM (
    'family_and_friend',
    'pre_seed',
    'bridge_round',
    'seed',
    'series_a',
    'series_b',
    'series_c'
);

CREATE TYPE fundraise_pipeline_column_type AS ENUM (
    'normal',
    'closed'
);

CREATE TABLE fundraising_pipelines (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    stage fundraise_pipeline_stage NOT NULL,
    reference VARCHAR(220) UNIQUE NOT NULL,
    workspace_id uuid NOT NULL REFERENCES workspaces(id),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    is_closed BOOLEAN NOT NULL DEFAULT false,
    target_amount BIGINT NOT NULL DEFAULT 0,
    closed_amount BIGINT NOT NULL DEFAULT 0,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expected_close_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE fundraising_pipelines ADD CONSTRAINT fundraising_pipeline_reference_check_key 
    CHECK (reference ~ 'fundraising_pipeline_[a-zA-Z0-9._]+');

CREATE TABLE fundraising_pipeline_columns (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    reference VARCHAR(220) UNIQUE NOT NULL,
    fundraising_pipeline_id uuid NOT NULL REFERENCES fundraising_pipelines(id),
    title TEXT NOT NULL,
    column_type fundraise_pipeline_column_type NOT NULL,
    description TEXT NOT NULL,
    investors_count BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE fundraising_pipeline_columns ADD CONSTRAINT fundraising_pipeline_column_reference_check_key 
    CHECK (reference ~ 'fundraising_pipeline_column_[a-zA-Z0-9._]+');

-- Create trigger to update updated_at on fundraising_pipelines
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_fundraising_pipelines_updated_at
    BEFORE UPDATE ON fundraising_pipelines
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_fundraising_pipeline_columns_updated_at
    BEFORE UPDATE ON fundraising_pipeline_columns
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
