BEGIN;

SET search_path TO uploader_service;

ALTER TABLE upload_events ADD COLUMN IF NOT EXISTS multipart_upload_id TEXT;
ALTER TABLE upload_events ADD COLUMN IF NOT EXISTS is_multipart BOOLEAN;

COMMIT;
