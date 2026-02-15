package usecase

import (
	"context"
	"time"

	"github.com/arunshankar19/home-server-common-utils/logger"
	"github.com/arunshankar19/home-server-common-utils/storage"
	"github.com/arunshankar19/home-server-upload-service/domain"
)

const (
	failedStatus = "failed"
)

type uploadBatchUsecase struct {
	log                   logger.Logger
	bucketName            string
	storage               storage.Storage
	uploadBatchRepository domain.UploadBatchRepository
}

// NewUploadBatchUsecase returns an implementation of batch usecase
func NewUploadBatchUsecase(
	log logger.Logger,
	bucketName string,
	storage storage.Storage,
	uploadBatchRepo domain.UploadBatchRepository,
) domain.UploadBatchUsecase {
	return &uploadBatchUsecase{
		log:                   log,
		bucketName:            bucketName,
		storage:               storage,
		uploadBatchRepository: uploadBatchRepo,
	}
}

// AbortStaleMultipartUploads removes stale uploads from storage and marks as failed
func (u *uploadBatchUsecase) AbortStaleMultipartUploads(
	ctx context.Context,
	batchSize int,
	uploadExpiryDuration time.Duration,
) error {
	var (
		offset   = 0
		fromTime = time.Now().Add(-uploadExpiryDuration)

		// batch stats
		noOfRecordsProcessed = 0
		batchStartTime       = time.Now()
	)

	for {
		staleUploads, err := u.uploadBatchRepository.FetchStaleUploads(ctx, fromTime, batchSize, offset)
		if err != nil {
			u.log.Error("failed to fetch stale uploads from db", map[string]any{"error": err})
			return err
		}

		if len(staleUploads) == 0 {
			break
		}

		successfullAborts := make([]string, 0, len(staleUploads))

		for _, staleUpload := range staleUploads {
			err = u.storage.AbortMultipartUpload(ctx, u.bucketName, staleUpload.ObjectName, staleUpload.MultipartUploadID)
			if err != nil {
				u.log.Error("failed to abort multipart from storage", map[string]any{"error": err})
				continue
			}
			successfullAborts = append(successfullAborts, staleUpload.ID)
		}

		err = u.uploadBatchRepository.UpdateStaleUploadStatus(ctx, failedStatus, successfullAborts)
		if err != nil {
			u.log.Error("failed to update failed upload status", map[string]any{"error": err})
		}

		offset += batchSize
		noOfRecordsProcessed += len(staleUploads)
	}

	u.log.Info("completed stale upload cleanup job", map[string]any{
		"no_of_records_processed":     noOfRecordsProcessed,
		"batch_start_time":            batchStartTime,
		"batch_end_time":              time.Now(),
		"batch_time_taken_in_seconds": time.Since(batchStartTime).Seconds(),
	})

	return nil
}
