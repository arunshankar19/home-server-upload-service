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
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id;
	`

	queryUpdateUploadEventStatus = `
	UPDATE upload_events SET status = $1 WHERE id = $2;
	`

	queryFindUploadOwner = `
	SELECT created_by FROM upload_events WHERE id = $1;
	`
)

// NewUploadRepository returns an implementation of upload repository
func NewUploadRepository(db database.DB, tracer observability.Tracer) domain.UploadRepository {
	return &uploadRepository{db: db, tracer: tracer}
}

// InsertUploadEvent inserts an upload event to db
func (r *uploadRepository) InsertUploadEvent(ctx context.Context, uploadEvent domain.UploadEvent) (string, error) {
	ctx, span := r.tracer.StartSpan(ctx, "Repository-InsertUploadEvent")
	defer span.End()

	var uploadEventID string
	err := r.db.QueryRow(ctx, queryInsertUploadEvents,
		uploadEvent.FileName,
		uploadEvent.FileType,
		uploadEvent.FileExt,
		uploadEvent.FileSize,
		uploadEvent.Status,
		uploadEvent.CreatedBy,
		time.Now(),
	).Scan(&uploadEventID)
	if err != nil {
		return "", err
	}
	return uploadEventID, nil
}

// UpdateUploadEventStatus updates the status of a given upload event
func (r *uploadRepository) UpdateUploadEventStatus(ctx context.Context, uploadEventID string, status string) error {
	ctx, span := r.tracer.StartSpan(ctx, "Repository-UpdateUploadEventStatus")
	defer span.End()

	_, err := r.db.Exec(ctx, queryUpdateUploadEventStatus, status, uploadEventID)
	return err
}

// FindUploadOwner retuns the id of uploaded user
func (r *uploadRepository) FindUploadOwner(ctx context.Context, uploadEventID string) (string, error) {
	var userID string

	err := r.db.QueryRow(ctx, queryFindUploadOwner, uploadEventID).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
