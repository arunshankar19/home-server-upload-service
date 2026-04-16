// this is the batch processing job that runs periodically
package main

import (
	"context"
	"log"
	"time"

	"github.com/arunshankar19/home-server-common-utils/logger"
	"github.com/arunshankar19/home-server-common-utils/storage"
	"github.com/arunshankar19/home-server-upload-service/internal/config"
	"github.com/arunshankar19/home-server-upload-service/internal/database"
	uploadRepository "github.com/arunshankar19/home-server-upload-service/modules/upload/repository/postgres"
	uploadUsecase "github.com/arunshankar19/home-server-upload-service/modules/upload/usecase"
)

func main() {
	ctx := context.Background()

	l, err := logger.NewLogger()
	if err != nil {
		log.Println("failed to initialize logger with error", err)
		panic(err)
	}
	appConfig, err := config.NewAppConfig()
	if err != nil {
		l.Error("failed to get all required envs", map[string]any{
			"error": err,
		})
		panic(err)
	}

	dbHost := database.WithHost(appConfig.PostgresHost)
	dbUser := database.WithUser(appConfig.PostgresUser)
	dbPassword := database.WithPassword(appConfig.PostgresPassword)
	dbName := database.WithDBName(appConfig.PostgresDBName)
	dbSearchPath := database.WithSearchPath(appConfig.PostgresSearchPath)
	dbPort := database.WithPort(appConfig.PostgresPort)
	dbSSLMode := database.WithSslMode(appConfig.PostgresSSLMode)

	db, err := database.NewPostgresDB(
		ctx,
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbSearchPath,
		dbPort,
		dbSSLMode,
	)
	if err != nil {
		l.Error("failed to connect to db", map[string]any{"error": err})
		panic(err)
	}
	err = db.Ping(ctx)
	if err != nil {
		l.Error("failed to ping db", map[string]any{"error": err})
		panic(err)
	}

	storageAddr := storage.WithEndpoint(appConfig.MinioAddr)
	storageAccesskey := storage.WithAccessKeyID(appConfig.MinioAccessKeyID)
	secretAccessKey := storage.WithSecretAccessKey(appConfig.MinioSecretAccessKey)
	storage, err := storage.NewMinioStorageClient(storageAddr, storageAccesskey, secretAccessKey)
	if err != nil {
		l.Error("failed to initialize storage client", map[string]any{"error": err})
		panic(err)
	}

	uploadBatchRepository := uploadRepository.NewUploadBatchRepo(db)
	uploadBatchUsecase := uploadUsecase.NewUploadBatchUsecase(
		l,
		appConfig.MinioDataBucketName,
		storage,
		uploadBatchRepository,
	)

	err = uploadBatchUsecase.AbortStaleMultipartUploads(ctx, appConfig.BatchSize, time.Duration(appConfig.UploadExpiryDurationInHours)*time.Hour)
	if err != nil {
		l.Error("failed to abort stale multipart uploads", map[string]any{"error": err})
		panic(err)
	}

	// log batch job completion successfully
	l.Info("completed batch job to abort stale multipart uploads", nil)
}
