package main

import (
	"context"
	"log"
	"net/http"

	"github.com/arunshankar19/home-server-common-utils/logger"
	"github.com/arunshankar19/home-server-common-utils/observability"
	"github.com/arunshankar19/home-server-common-utils/storage"
	"github.com/arunshankar19/home-server-upload-service/domain"
	"github.com/arunshankar19/home-server-upload-service/internal/config"
	"github.com/arunshankar19/home-server-upload-service/internal/database"
	"github.com/arunshankar19/home-server-upload-service/internal/middlewares"
	uploadHandler "github.com/arunshankar19/home-server-upload-service/modules/upload/delivery/http"
	uploadRepository "github.com/arunshankar19/home-server-upload-service/modules/upload/repository/postgres"
	uploadUsecase "github.com/arunshankar19/home-server-upload-service/modules/upload/usecase"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	serviceName = "upload-service"
)

type handlers struct {
	uploadHandler domain.UploadHandler
}

func main() {
	l, err := logger.NewLogger()
	if err != nil {
		log.Println("failed to initialize logger with error", err)
		panic(err)
	}

	ctx := context.Background()

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	appConfig, err := config.NewAppConfig()
	if err != nil {
		l.Error("failed to get all required envs", map[string]any{
			"error": err,
		})
		panic(err)
	}

	l.Info("application configs read successfully", map[string]any{
		"appConfig": appConfig,
	})

	// initialise tracer
	trace := observability.NewNoopTracer()
	if appConfig.ObservabilityEnabled {
		observability.InitTraceProvider(ctx, appConfig.TraceExpoterURL, serviceName)
		trace = observability.NewTracer(serviceName)
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

	uploadRepo := uploadRepository.NewUploadRepository(db, trace)
	uploadUsecae := uploadUsecase.NewUploadUsecase(l, appConfig.MinioDataBucketName, appConfig.MinioPresignedURLExpInMinutes, storage, uploadRepo)
	uploadHandler := uploadHandler.NewUploadHandler(uploadUsecae)

	httpHandler := handlers{
		uploadHandler: uploadHandler,
	}

	initV1Router(r, httpHandler)

	srv := http.Server{
		Addr:    ":8001",
		Handler: r,
	}

	srv.ListenAndServe()
}

func initV1Router(r *chi.Mux, h handlers) {
	apiV1Router := chi.NewRouter()

	apiV1Router.Use(middlewares.SetUserIdInCtx)

	apiV1Router.Get("/status", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte{})
	})
	apiV1Router.Post("/init-upload", h.uploadHandler.InitiateUpload)
	apiV1Router.Post("/get-presigned-url", h.uploadHandler.GetPresignedURL)
	apiV1Router.Post("/complete-upload", h.uploadHandler.CompleteUpload)

	r.Mount("/v1", apiV1Router)
}
