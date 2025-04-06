CREATE TYPE fundraising_column_activity AS ENUM (
    'meeting',
    'note',
    'email'
);

CREATE TABLE fundraising_pipeline_column_contacts (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    reference VARCHAR(220) UNIQUE NOT NULL,
    contact_id uuid NOT NULL REFERENCES contacts(id),
    fundraising_pipeline_id uuid NOT NULL REFERENCES fundraising_pipelines(id),
    fundraising_pipeline_column_id uuid NOT NULL REFERENCES fundraising_pipeline_columns(id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(fundraising_pipeline_id, contact_id)
);

ALTER TABLE fundraising_pipeline_column_contacts ADD CONSTRAINT fundraising_pipeline_column_contact_reference_check_key 
    CHECK (reference ~ 'fundraising_pipeline_column_contact_[a-zA-Z0-9._]+');

CREATE TABLE fundraising_pipeline_column_contact_positions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    reference VARCHAR(220) UNIQUE NOT NULL,
    fundraising_pipeline_column_contact_id uuid NOT NULL REFERENCES fundraising_pipeline_column_contacts(id),
    order_index BIGINT NOT NULL,
    UNIQUE(fundraising_pipeline_column_contact_id)
);

ALTER TABLE fundraising_pipeline_column_contact_positions ADD CONSTRAINT fundraising_pipeline_column_contact_position_reference_check_key 
    CHECK (reference ~ 'fundraising_pipeline_column_contact_position_[a-zA-Z0-9._]+');

CREATE TABLE fundraising_pipeline_column_contact_deals (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    reference VARCHAR(220) UNIQUE NOT NULL,
    fundraising_pipeline_column_contact_id uuid NOT NULL REFERENCES fundraising_pipeline_column_contacts(id),
    check_size BIGINT NOT NULL DEFAULT 0,
    can_lead_round BOOLEAN NOT NULL DEFAULT false,
    rating BIGINT NOT NULL DEFAULT 0,
    initial_contact TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(fundraising_pipeline_column_contact_id)
);

ALTER TABLE fundraising_pipeline_column_contact_deals ADD CONSTRAINT fundraising_pipeline_column_contact_deal_reference_check_key 
    CHECK (reference ~ 'fundraising_pipeline_column_contact_deal_[a-zA-Z0-9._]+');

CREATE TABLE fundraising_pipeline_column_contact_activities (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    reference VARCHAR(220) UNIQUE NOT NULL,
    fundraising_pipeline_column_contact_id uuid NOT NULL REFERENCES fundraising_pipeline_column_contacts(id),
    activity_type fundraising_column_activity NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE fundraising_pipeline_column_contact_activities ADD CONSTRAINT fundraising_pipeline_column_contact_activity_reference_check_key 
    CHECK (reference ~ 'fundraising_pipeline_column_contact_activity_[a-zA-Z0-9._]+');

CREATE TABLE fundraising_pipeline_column_contact_documents (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    reference VARCHAR(220) UNIQUE NOT NULL,
    fundraising_pipeline_column_contact_id uuid NOT NULL REFERENCES fundraising_pipeline_column_contacts(id),
    title TEXT NOT NULL,
    file_size BIGINT NOT NULL DEFAULT 0,
    object_key TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE fundraising_pipeline_column_contact_documents ADD CONSTRAINT fundraising_pipeline_column_contact_document_reference_check_key 
    CHECK (reference ~ 'fundraising_pipeline_column_contact_document_[a-zA-Z0-9._]+');

CREATE TRIGGER update_fundraising_pipeline_column_contacts_updated_at
    BEFORE UPDATE ON fundraising_pipeline_column_contacts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_fundraising_pipeline_column_contact_deals_updated_at
    BEFORE UPDATE ON fundraising_pipeline_column_contact_deals
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_fundraising_pipeline_column_contact_documents_updated_at
    BEFORE UPDATE ON fundraising_pipeline_column_contact_documents
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
