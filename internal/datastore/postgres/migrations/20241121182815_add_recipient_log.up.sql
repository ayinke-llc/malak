CREATE TYPE update_recipient_email_provider AS ENUM ('resend', 'sendgrid','smtp');

CREATE TABLE update_recipient_logs(
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  reference VARCHAR (220) UNIQUE NOT NULL,
  recipient_id uuid NOT NULL REFERENCES update_recipients(id),
  provider_id VARCHAR (220) UNIQUE NOT NULL,
  provider update_recipient_email_provider NOT NULL,

  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
);

ALTER TABLE update_recipient_logs
  ADD CONSTRAINT updates_reference_logs_check_key 
  CHECK (reference ~ 'recipient_log_[a-zA-Z0-9._]+');
