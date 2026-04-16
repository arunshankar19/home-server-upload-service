package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/arunshankar19/home-server-upload-service/domain"
	"github.com/arunshankar19/home-server-upload-service/internal/response"
)

type uploadHandler struct {
	uploadUsecase domain.UploadUsecase
}

// NewUploadHandler returns a upload handler
func NewUploadHandler(uploadUsecase domain.UploadUsecase) domain.UploadHandler {
	return &uploadHandler{
		uploadUsecase: uploadUsecase,
	}
}

// InitiateUpload starts an multipart upload
func (h *uploadHandler) InitiateUpload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var initiateUploadReq domain.InitiateUploadReqDTO

	err := json.NewDecoder(r.Body).Decode(&initiateUploadReq)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	initiateUploadResult, err := h.uploadUsecase.InitiateUpload(ctx, initiateUploadReq)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedUser) {
			response.Fail(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			return
		}
		response.Fail(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(w, http.StatusOK, domain.InitiateUploadResponseDTO{
		UploadID:      initiateUploadResult.UploadID,
		UploadEventID: initiateUploadResult.UploadEventID,
	})
}

// GetPresignedURL handles presigned url http handler
func (h *uploadHandler) GetPresignedURL(w http.ResponseWriter, r *http.Request) {
	var presignedURLReq domain.PresignedMultipartURLReqDTO

	err := json.NewDecoder(r.Body).Decode(&presignedURLReq)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	url, err := h.uploadUsecase.GetPresignedURL(r.Context(), presignedURLReq)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(w, http.StatusOK, domain.PresignedMultipartResponseDTO{PresignedURL: url})
}

// CompleteUpload is the handler for committing the multipart upload
func (h *uploadHandler) CompleteUpload(w http.ResponseWriter, r *http.Request) {
	var completeUploadReq domain.CompleteUploadReqDTO

	err := json.NewDecoder(r.Body).Decode(&completeUploadReq)
	if err != nil {
		response.Fail(w, http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	err = h.uploadUsecase.CompleteUpload(r.Context(), completeUploadReq)
	if err != nil {
		response.Fail(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	response.Success(w, http.StatusOK, nil)
}
