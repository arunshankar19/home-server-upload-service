package repository

import (
	"context"
	"time"

	"github.com/arunshankar19/home-server-common-utils/observability"
	"github.com/arunshankar19/home-server-upload-service/domain"
	"github.com/arunshankar19/home-server-upload-service/internal/database"
)

type uploadRepository struct {
	db     database.DB
	tracer observability.Tracer
}

// queries
const (
	queryInsertUploadEvents = `
	INSERT INTO upload_events (
		file_name,
		file_type,
		file_ext,
		file_size,
		status,
		created_by,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7);
	`
)

// NewUploadRepository returns an implementation of upload repository
func NewUploadRepository(db database.DB, tracer observability.Tracer) domain.UploadRepository {
	return &uploadRepository{db: db, tracer: tracer}
}

// InsertUploadEvent inserts an upload event to db
func (r *uploadRepository) InsertUploadEvent(ctx context.Context, uploadEvent domain.UploadEvent) error {
	ctx, span := r.tracer.StartSpan(ctx, "Repository-InsertUploadEvent")
	defer span.End()

	_, err := r.db.Exec(ctx, queryInsertUploadEvents,
		uploadEvent.FileName,
		uploadEvent.FileType,
		uploadEvent.FileExt,
		uploadEvent.FileSize,
		uploadEvent.Status,
		uploadEvent.CreatedBy,
		time.Now(),
	)
	return err
}
