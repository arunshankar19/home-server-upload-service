BEGIN;

SET search_path TO uploader_service;

ALTER TABLE upload_events DROP COLUMN IF EXISTS multipart_upload_id;
ALTER TABLE upload_events DROP COLUMN IF EXISTS is_multipart;

COMMIT;
