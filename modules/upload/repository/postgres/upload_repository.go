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
		id,
		file_name,
		file_type,
		file_ext,
		file_size,
		is_multipart,
		multipart_upload_id,
		status,
		created_by,
		created_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	queryUpdateUploadEventStatus = `
	UPDATE upload_events SET status = $1 WHERE id = $2;
	`

	queryFindUploadOwner = `
	SELECT created_by FROM upload_events WHERE id = $1;
	`

	queryFetchUploadEvent = `
	SELECT
		id,
		file_name,
		file_type,
		file_ext,
		file_size,
		is_multipart,
		multipart_upload_id,
		status,
		created_by,
		created_at
	FROM
		upload_events
	WHERE
		id = $1;
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

	var (
		multipartUploadID *string
	)

	if uploadEvent.IsMultipart {
		multipartUploadID = &uploadEvent.MultipartUploadID
	}

	_, err := r.db.Exec(ctx, queryInsertUploadEvents,
		uploadEvent.ID,
		uploadEvent.FileName,
		uploadEvent.FileType,
		uploadEvent.FileExt,
		uploadEvent.FileSize,
		uploadEvent.IsMultipart,
		multipartUploadID,
		uploadEvent.Status,
		uploadEvent.CreatedBy,
		time.Now(),
	)
	return err
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
	ctx, span := r.tracer.StartSpan(ctx, "Repository-FindUploadOwner")
	defer span.End()

	var userID string

	err := r.db.QueryRow(ctx, queryFindUploadOwner, uploadEventID).Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}

// GetUploadEvent retrieves upload event based on id
func (r *uploadRepository) GetUploadEvent(ctx context.Context, uploadEventID string) (*domain.UploadEvent, error) {
	ctx, span := r.tracer.StartSpan(ctx, "Repository-GetUploadEvent")
	defer span.End()

	var uploadEvent domain.UploadEvent
	err := r.db.QueryRow(ctx, queryFetchUploadEvent, uploadEventID).Scan(
		&uploadEvent.ID,
		&uploadEvent.FileName,
		&uploadEvent.FileType,
		&uploadEvent.FileExt,
		&uploadEvent.FileSize,
		&uploadEvent.IsMultipart,
		&uploadEvent.MultipartUploadID,
		&uploadEvent.Status,
		&uploadEvent.CreatedBy,
		&uploadEvent.CreatedAt,
	)
	if err != nil {
		if database.IsNullErr(err) {
			return nil, err
		}
		return nil, err
	}
	return &uploadEvent, nil
}
