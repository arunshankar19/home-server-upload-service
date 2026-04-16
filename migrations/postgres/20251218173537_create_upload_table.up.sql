BEGIN;

SET search_path TO uploader_service;

CREATE TABLE IF NOT EXISTS upload_events (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_name   TEXT        NOT NULL,
    file_type   TEXT        NOT NULL,
    file_ext    TEXT        NOT NULL,
    file_size   BIGINT      NOT NULL,
    status      TEXT        NOT NULL,
    created_by  UUID        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);


COMMIT;