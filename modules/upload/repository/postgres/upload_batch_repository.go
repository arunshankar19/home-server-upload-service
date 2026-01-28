package repository

import (
	"context"
	"time"

	"github.com/arunshankar19/home-server-upload-service/domain"
	"github.com/arunshankar19/home-server-upload-service/internal/database"
)

type uploadBatchRepo struct {
	db database.DB
}

const (
	queryUpdateStateUploadStatus = `
	UPDATE upload_events SET status = $1 WHERE id = ANY($2);
	`

	queryFetchStaleUploads = `
	SELECT
		id,
		multipart_upload_id,
		file_name
	FROM
		upload_events
	WHERE
		status = 'pending'
	AND
		updated_at <= $1
	ORDER BY updated_at
	OFFSET $2
	LIMIT $3;
	`
)

// NewUploadBatchRepo returns a upload batch repository
func NewUploadBatchRepo(db database.DB) domain.UploadBatchRepository {
	return &uploadBatchRepo{db: db}
}

// UpdateStaleUploadStatus updates status which is before fromTime
func (r *uploadBatchRepo) UpdateStaleUploadStatus(
	ctx context.Context,
	status string,
	ids []string,
) error {
	_, err := r.db.Exec(ctx, queryUpdateStateUploadStatus, status, ids)
	return err
}

// FetchStaleUploads fetches stale uploads details
func (r *uploadBatchRepo) FetchStaleUploads(
	ctx context.Context,
	fromTime time.Time,
	limit int,
	offset int,
) ([]domain.StaleMultipartUpload, error) {
	staleUploads := make([]domain.StaleMultipartUpload, 0, limit)

	rows, err := r.db.Query(ctx, queryFetchStaleUploads, fromTime, offset, limit)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var staleUpload domain.StaleMultipartUpload
		err = rows.Scan(
			&staleUpload.ID,
			&staleUpload.MultipartUploadID,
			&staleUpload.ObjectName,
		)
		if err != nil {
			return nil, err
		}
		staleUploads = append(staleUploads, staleUpload)
	}

	return staleUploads, nil
}
