CREATE TABLE pipelines (
  id UUID PRIMARY KEY,
  project_id UUID REFERENCES projects(id),
  stage pipeline_stage NOT NULL DEFAULT 'initial',
  stage_status stage_status NOT NULL DEFAULT 'pending',
  fail_description VARCHAR NOT NULL,
  input_url VARCHAR NOT NULL,
  output_key VARCHAR NOT NULL,
  output_path VARCHAR NOT NULL,
  size_kb FLOAT NOT NULL,
  transcode_duration_seconds FLOAT NOT NULL,
  upload_duration_seconds FLOAT NOT NULL,
  max_resolution resolutions_enum NOT NULL,
  resolutions resolutions_enum[] NOT NULL DEFAULT '{}',
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
