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

// InitUploadResult holds data from init upload
type InitUploadResult struct {
	UploadID      string
	UploadEventID string
}

// UploadRepository defines the repository funcs
type UploadRepository interface {
	InsertUploadEvent(ctx context.Context, uploadEvent UploadEvent) (string, error)
	UpdateUploadEventStatus(ctx context.Context, uploadEventID string, status string) error
	FindUploadOwner(ctx context.Context, uploadEventID string) (string, error)
}

// UploadUsecase defines business usecase funcs
type UploadUsecase interface {
	InitiateMultipartUpload(ctx context.Context, initiateUploadReq InitiateUploadReqDTO) (*InitUploadResult, error)
	GetPresignedURL(ctx context.Context, presignedURLReq PresignedMultipartURLReqDTO) (string, error)
	CompleteUpload(ctx context.Context, completeUploadReq CompleteUploadReqDTO) error
}

// UploadHandler defines http handlers
type UploadHandler interface {
	InitiateMultipartUpload(w http.ResponseWriter, r *http.Request)
	GetPresignedURL(w http.ResponseWriter, r *http.Request)
	CompleteMultipartUpload(w http.ResponseWriter, r *http.Request)
}
