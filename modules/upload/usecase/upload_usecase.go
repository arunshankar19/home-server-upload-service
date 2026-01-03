package usecase

import (
	"context"

	"github.com/arunshankar19/home-server-common-utils/logger"
	"github.com/arunshankar19/home-server-common-utils/storage"
	"github.com/arunshankar19/home-server-upload-service/domain"
	"github.com/arunshankar19/home-server-upload-service/internal/ctxkey"
	"github.com/google/uuid"
)

const (
	uploadPending = "pending"
)

type uploadUsecase struct {
	log              logger.Logger
	uploadBucket     string
	storage          storage.Storage
	uploadRepository domain.UploadRepository
}

// NewUploadUsecase returns a upload usecase implementation
func NewUploadUsecase(
	log logger.Logger,
	uploadBucket string,
	storage storage.Storage,
	uploadRepository domain.UploadRepository,
) domain.UploadUsecase {
	return &uploadUsecase{
		log:              log,
		uploadBucket:     uploadBucket,
		storage:          storage,
		uploadRepository: uploadRepository,
	}
}

// InitiateMultipartUpload starts a multipart upload and marks it in db
func (u *uploadUsecase) InitiateMultipartUpload(
	ctx context.Context,
	initiateUploadReq domain.InitiateUploadReqDTO,
) (string, error) {

	userIDString, ok := ctx.Value(ctxkey.UserID).(string)
	if !ok {
		u.log.Error("unauthorised user - there is no user id in context", nil)
		return "", domain.ErrUnauthorizedUser
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		u.log.Error("failed to initialise multipart upload", map[string]any{"error": err})
		return "", err
	}

	uploadID, err := u.storage.InitialiseMultipartUpload(ctx, u.uploadBucket, initiateUploadReq.FileName, storage.PutObjectOptions{ContentType: initiateUploadReq.FileType})
	if err != nil {
		u.log.Error("failed to initialise multipart upload", map[string]any{"error": err})
		return "", err
	}

	uploadEvent := domain.UploadEvent{
		FileName:  initiateUploadReq.FileName,
		FileType:  initiateUploadReq.FileType,
		FileExt:   initiateUploadReq.FileExtension,
		FileSize:  initiateUploadReq.FileSize,
		Status:    uploadPending,
		CreatedBy: userID,
	}

	err = u.uploadRepository.InsertUploadEvent(ctx, uploadEvent)
	if err != nil {
		return "", err
	}
	return uploadID, nil
}
