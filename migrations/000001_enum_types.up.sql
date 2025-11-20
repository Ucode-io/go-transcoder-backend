-- Create the role_enum enumeration type
CREATE TYPE role_enum AS ENUM (
  'superadmin',
  'project'
);

-- Create the status_enum enumeration type
CREATE TYPE status_enum AS ENUM (
  'active',
  'inactive',
  'deleted'
);

-- Create the pipeline_stage enumeration type
CREATE TYPE pipeline_stage AS ENUM (
  'initial',
  'preparation',
  'transcode',
  'upload'
);

-- Create the pipeline_stage_status enumeration type
CREATE TYPE stage_status AS ENUM (
  'pending',
  'success',
  'fail'
);

-- Create the resolutions_enum enumeration type
CREATE TYPE resolutions_enum AS ENUM (
  '240p',
  '360p',
  '480p',
  '720p',
  '1080p',
  '1440p',
  '4k'
);

-- Create the storage_type enumeration type
CREATE TYPE storage_type AS ENUM (
  'minio',
  's3',
  'unknown'
);

-- Create the webhook_type enumeration type
CREATE TYPE webhook_type AS ENUM (
  'update_status'
);
