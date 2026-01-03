package domain

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// UploadEvent represents an upload event in db
type UploadEvent struct {
	ID        uuid.UUID
	FileName  string
	FileType  string
	FileExt   string
	FileSize  int
	Status    string
	CreatedBy uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

// UploadRepository defines the repository funcs
type UploadRepository interface {
	InsertUploadEvent(ctx context.Context, uploadEvent UploadEvent) error
}

// UploadUsecase defines business usecase funcs
type UploadUsecase interface {
	InitiateMultipartUpload(ctx context.Context, initiateUploadReq InitiateUploadReqDTO) (string, error)
}

// UploadHandler defines http handlers
type UploadHandler interface {
	InitiateMultipartUpload(w http.ResponseWriter, r *http.Request)
}
