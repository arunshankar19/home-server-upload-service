package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/arunshankar19/home-server-common-utils/logger"
	"github.com/arunshankar19/home-server-common-utils/observability"
	"github.com/arunshankar19/home-server-common-utils/storage"
	"github.com/arunshankar19/home-server-upload-service/domain"
	"github.com/arunshankar19/home-server-upload-service/internal/ctxkey"
	fileexplorer "github.com/arunshankar19/home-server-upload-service/internal/file_explorer"
	"github.com/google/uuid"
)

const (
	uploadPending   = "pending"
	uploadCompleted = "completed"
)

type uploadUsecase struct {
	log                logger.Logger
	tracer             observability.Tracer
	uploadBucket       string
	preSignedURLExp    time.Duration
	storage            storage.Storage
	fileExplorerClient fileexplorer.FileExplorerClient
	uploadRepository   domain.UploadRepository
}

// NewUploadUsecase returns a upload usecase implementation
func NewUploadUsecase(
	log logger.Logger,
	tracer observability.Tracer,
	uploadBucket string,
	presignedURLExp int,
	storage storage.Storage,
	fileExplorerClient fileexplorer.FileExplorerClient,
	uploadRepository domain.UploadRepository,
) domain.UploadUsecase {
	return &uploadUsecase{
		log:                log,
		tracer:             tracer,
		uploadBucket:       uploadBucket,
		preSignedURLExp:    time.Duration(presignedURLExp) * time.Minute,
		storage:            storage,
		fileExplorerClient: fileExplorerClient,
		uploadRepository:   uploadRepository,
	}
}

// InitiateUpload starts a multipart upload and marks it in db
func (u *uploadUsecase) InitiateUpload(
	ctx context.Context,
	initiateUploadReq domain.InitiateUploadReqDTO,
) (*domain.InitUploadResult, error) {

	ctx, span := u.tracer.StartSpan(ctx, "Usecase-InitiateUpload")
	defer span.End()

	userIDString, ok := ctx.Value(ctxkey.UserID).(string)
	if !ok {
		u.log.Error("unauthorised user - there is no user id in context", nil)
		return nil, domain.ErrUnauthorizedUser
	}

	userID, err := uuid.Parse(userIDString)
	if err != nil {
		u.log.Error("invalid user id", map[string]any{"error": err})
		return nil, domain.ErrInvalidUserID
	}

	// this will be the upload event id and file name in storage
	uploadEventID := uuid.New()

	var uploadID string
	// initiate a multipart session if a multipart upload
	if initiateUploadReq.IsMultipart {
		uploadID, err = u.storage.InitialiseMultipartUpload(
			ctx,
			u.uploadBucket,
			fmt.Sprintf("%s.%s", uploadEventID.String(), initiateUploadReq.FileExtension),
			storage.PutObjectOptions{ContentType: initiateUploadReq.FileType},
		)
		if err != nil {
			u.log.Error("failed to initialise multipart upload", map[string]any{"error": err})
			return nil, err
		}
	}

	uploadEvent := domain.UploadEvent{
		ID:                uploadEventID,
		FileName:          initiateUploadReq.FileName,
		FileType:          initiateUploadReq.FileType,
		FileExt:           initiateUploadReq.FileExtension,
		FileSize:          initiateUploadReq.FileSize,
		IsMultipart:       initiateUploadReq.IsMultipart,
		MultipartUploadID: uploadID,
		Status:            uploadPending,
		CreatedBy:         userID,
	}

	err = u.uploadRepository.InsertUploadEvent(ctx, uploadEvent)
	if err != nil {
		u.log.Error("failed to insert event in db", map[string]any{"error": err})
		return nil, err
	}
	return &domain.InitUploadResult{UploadID: uploadID, UploadEventID: uploadEventID.String()}, nil
}

// GetPresignedURL returnes a presigned url with given expiry
func (u *uploadUsecase) GetPresignedURL(
	ctx context.Context,
	presignedURLReq domain.PresignedMultipartURLReqDTO,
) (string, error) {

	ctx, span := u.tracer.StartSpan(ctx, "Usecase-GetPresignedURL")
	defer span.End()

	userIDString, ok := ctx.Value(ctxkey.UserID).(string)
	if !ok {
		u.log.Error("unauthorised user - there is no user id in context", nil)
		return "", domain.ErrUnauthorizedUser
	}

	_, err := uuid.Parse(userIDString)
	if err != nil {
		u.log.Error("failed to initialise multipart upload", map[string]any{"error": err})
		return "", domain.ErrInvalidUserID
	}

	uploadEvent, err := u.uploadRepository.GetUploadEvent(ctx, presignedURLReq.UploadEventID)
	if err != nil {
		u.log.Error("failed to find file upload owner", map[string]any{"error": err})
		return "", err
	}

	if userIDString != uploadEvent.CreatedBy.String() {
		u.log.Error("the user is not the owner of the file", nil)
		return "", domain.ErrNotFileOwner
	}

	var reqParams map[string][]string
	if presignedURLReq.IsMultipart {
		reqParams = make(map[string][]string)
		reqParams["partNumber"] = []string{fmt.Sprintf("%d", presignedURLReq.PartNumber)}
		reqParams["uploadId"] = []string{uploadEvent.MultipartUploadID}
	}

	url, err := u.storage.GetPresignedURL(
		ctx,
		u.uploadBucket,
		fmt.Sprintf("%s.%s", presignedURLReq.UploadEventID, uploadEvent.FileExt),
		u.preSignedURLExp,
		reqParams,
	)
	if err != nil {
		u.log.Error("failed to generate presigned url", map[string]any{"error": err})
		return "", err
	}

	return url.String(), nil
}

// CompleteUpload updates the db status to complete
// if multipart upload it commits the upload
func (u *uploadUsecase) CompleteUpload(
	ctx context.Context,
	completeUploadReq domain.CompleteUploadReqDTO,
) error {

	ctx, span := u.tracer.StartSpan(ctx, "Usecase-CompleteUpload")
	defer span.End()

	userIDString, ok := ctx.Value(ctxkey.UserID).(string)
	if !ok {
		u.log.Error("unauthorised user - there is no user id in context", nil)
		return domain.ErrUnauthorizedUser
	}

	_, err := uuid.Parse(userIDString)
	if err != nil {
		u.log.Error("failed to initialise multipart upload", map[string]any{"error": err})
		return domain.ErrInvalidUserID
	}

	uploadEvent, err := u.uploadRepository.GetUploadEvent(ctx, completeUploadReq.UploadEventID)
	if err != nil {
		u.log.Error("failed to find file upload owner", map[string]any{"error": err})
		return err
	}
	if uploadEvent == nil {
		u.log.Error("there is no upload event for the given id", nil)
		return domain.ErrUploadEventNotPresent
	}

	if userIDString != uploadEvent.CreatedBy.String() {
		u.log.Error("the user is not the owner of the file", nil)
		return domain.ErrNotFileOwner
	}

	// commit the upload if multipart upload
	if uploadEvent.IsMultipart {
		completedParts := make([]storage.CompleteParts, 0, len(completeUploadReq.Etags))
		for _, partNumberEtagMapping := range completeUploadReq.Etags {
			completedParts = append(completedParts, storage.CompleteParts{
				PartNumber: partNumberEtagMapping.PartNumber,
				ETag:       partNumberEtagMapping.Etag,
			})
		}
		err = u.storage.CompleteMultipartUpload(
			ctx,
			u.uploadBucket,
			fmt.Sprintf("%s.%s", completeUploadReq.UploadEventID, uploadEvent.FileExt),
			uploadEvent.MultipartUploadID,
			completedParts,
			storage.PutObjectOptions{
				UserMetadata: map[string]string{"userID": userIDString},
			},
		)
		if err != nil {
			u.log.Error("failed to complete multipart upload", map[string]any{"error": err})
			return err
		}
	}

	err = u.fileExplorerClient.NotifyFileExplorer(ctx, *uploadEvent)
	if err != nil {
		u.log.Error("failed to notify file explorer service", map[string]any{"error": err})
		return err
	}

	err = u.uploadRepository.UpdateUploadEventStatus(ctx, completeUploadReq.UploadEventID, uploadCompleted)
	if err != nil {
		u.log.Error("failed to update upload events table", map[string]any{"error": err})
		return err
	}

	return nil
}
